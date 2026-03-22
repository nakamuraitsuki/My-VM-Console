package compute

import (
	"context"
	"errors"
	"fmt"
	"log"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/gateway"
	"example.com/m/internal/domain/image"
	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/storage"
	"example.com/m/internal/usecase"
)

type CreateInstancePayload struct {
	InstanceID compute.InstanceID `json:"instance_id"`
}

type ExecuteCreateInstanceUseCase interface {
	Execute(ctx context.Context, payload CreateInstancePayload) error
}

type executeCreateInstanceInteractor struct {
	instanceRepo   compute.InstanceRepository
	networkRepo    network.Repository
	ingressRepo    gateway.Repository
	storageRepo    storage.Repository
	imgRepo        image.Repository
	instanceDriver compute.ComputeDriver
	storageDriver  storage.StorageDriver
	gatewayDriver  gateway.IngressDriver
	uow            usecase.UnitOfWork
}

func NewExecuteCreateInstanceInteractor(
	instanceRepo compute.InstanceRepository,
	networkRepo network.Repository,
	storageRepo storage.Repository,
	ingressRepo gateway.Repository,
	imgRepo image.Repository,
	instanceDriver compute.ComputeDriver,
	storageDriver storage.StorageDriver,
	gatewayDriver gateway.IngressDriver,
	uow usecase.UnitOfWork,
) ExecuteCreateInstanceUseCase {
	return &executeCreateInstanceInteractor{
		instanceRepo:   instanceRepo,
		networkRepo:    networkRepo,
		storageRepo:    storageRepo,
		ingressRepo:    ingressRepo,
		imgRepo:        imgRepo,
		instanceDriver: instanceDriver,
		storageDriver:  storageDriver,
		gatewayDriver:  gatewayDriver,
		uow:            uow,
	}
}

func (i *executeCreateInstanceInteractor) Execute(ctx context.Context, payload CreateInstancePayload) error {
	var inst *compute.Instance
	var img *image.Image
	var vpc *network.VPC
	log.Printf("[ExecuteCreateInstance] start instanceID=%s", payload.InstanceID)

	if err := i.uow.Do(ctx, func(ctx context.Context) error {
		var uowErr error
		inst, uowErr = i.instanceRepo.FindByID(ctx, payload.InstanceID)
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: FindByID instanceID=%s err=%v", payload.InstanceID, uowErr)
			return compute.ErrInstanceNotFound
		}
		log.Printf("[ExecuteCreateInstance] instance loaded instanceID=%s status=%s", inst.ID(), inst.Status())

		img, uowErr = i.imgRepo.FindByID(ctx, inst.ImageID())
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: FindImage imageID=%s err=%v", inst.ImageID(), uowErr)
			return uowErr
		}
		if img == nil {
			log.Printf("[ExecuteCreateInstance] failed: image not found imageID=%s", inst.ImageID())
			return image.ErrImageNotFound
		}
		log.Printf("[ExecuteCreateInstance] image loaded imageID=%s name=%s", img.ID(), img.Alias())

		vpc, uowErr = i.networkRepo.FindVPCByUserID(ctx, inst.OwnerID())
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: FindVPC ownerID=%s err=%v", inst.OwnerID(), uowErr)
			return uowErr
		}
		if vpc == nil {
			log.Printf("[ExecuteCreateInstance] failed: VPC not found for ownerID=%s", inst.OwnerID())
			return fmt.Errorf("VPC not found for user: %s", inst.OwnerID())
		}
		log.Printf("[ExecuteCreateInstance] vpc loaded vpcID=%s ownerID=%s", vpc.ID(), vpc.OwnerID())

		uowErr = inst.MarkAsCreating() // 状態を「作成中」に遷移させる
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: MarkAsCreating instanceID=%s err=%v", inst.ID(), uowErr)
			return compute.ErrInvalidInstanceStatus
		}
		log.Printf("[ExecuteCreateInstance] marked creating instanceID=%s", inst.ID())

		return i.instanceRepo.Save(ctx, inst)
	}); err != nil {
		if errors.Is(err, compute.ErrInstanceNotFound) {
			log.Printf("[ExecuteCreateInstance] failed: uow (not found)")
			return compute.ErrInstanceNotFound
		}
		log.Printf("[ExecuteCreateInstance] failed: uow err=%v", err)
		inst.MarkAsError(compute.ErrInPending)
		return err
	}
	log.Printf("[ExecuteCreateInstance] uow committed (state=Creating) instanceID=%s", inst.ID())

	//　storage
	volumeID := inst.RootVolumeID()
	volume, err := i.storageRepo.FindByID(ctx, volumeID)
	if err != nil {
		log.Printf("[ExecuteCreateInstance] failed: FindVolume volumeID=%s err=%v", volumeID, err)
		return err
	}
	log.Printf("[ExecuteCreateInstance] volume loaded volumeID=%s name=%s", volume.ID(), volume.Name())
	log.Printf("[ExecuteCreateInstance] creating volume vpcID=%s volumeID=%s", vpc.ID(), volume.ID())
	err = i.storageDriver.CreateVolume(ctx, vpc.ID(), volume)
	if err != nil {
		log.Printf("[ExecuteCreateInstance] failed: CreateVolume volumeID=%s err=%v", volume.ID(), err)
		inst.MarkAsError(compute.ErrInCreating)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}
	log.Printf("[ExecuteCreateInstance] volume created volumeID=%s", volume.ID())

	// instance
	log.Printf("[ExecuteCreateInstance] creating instance driver instanceID=%s imageID=%s", inst.ID(), img.ID())
	err = i.instanceDriver.Create(ctx, inst, img)
	if err != nil {
		log.Printf("[ExecuteCreateInstance] failed: instanceDriver.Create instanceID=%s err=%v", inst.ID(), err)
		inst.MarkAsError(compute.ErrInCreating)
		_ = i.instanceRepo.Save(ctx, inst)
		_ = i.storageDriver.DeleteVolume(ctx, vpc.ID(), volume) // 作成したVolumeを削除
		log.Printf("[ExecuteCreateInstance] rollback: deleted volume volumeID=%s", volume.ID())
		return err
	}
	log.Printf("[ExecuteCreateInstance] instance created instanceID=%s", inst.ID())

	// 遷移
	log.Printf("[ExecuteCreateInstance] marking starting instanceID=%s", inst.ID())
	err = i.uow.Do(ctx, func(ctx context.Context) error {
		var uowErr error
		// 最新状態をDBから取得する
		inst, uowErr = i.instanceRepo.FindByID(ctx, payload.InstanceID)
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: FindByID (for starting) instanceID=%s err=%v", payload.InstanceID, uowErr)
			return compute.ErrInstanceNotFound
		}
		uowErr = inst.MarkAsStarting() // 状態を「起動中」に遷移させる
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: MarkAsStarting instanceID=%s err=%v", inst.ID(), uowErr)
			return compute.ErrInvalidInstanceStatus
		}
		log.Printf("[ExecuteCreateInstance] marked starting instanceID=%s", inst.ID())
		return i.instanceRepo.Save(ctx, inst)
	})
	if err != nil {
		log.Printf("[ExecuteCreateInstance] failed: uow (starting) err=%v", err)
		inst.MarkAsError(compute.ErrInCreating)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}
	log.Printf("[ExecuteCreateInstance] uow committed (state=Starting) instanceID=%s", inst.ID())

	log.Printf("[ExecuteCreateInstance] starting instance driver instanceID=%s", inst.ID())
	err = i.instanceDriver.Start(ctx, inst)
	if err != nil {
		log.Printf("[ExecuteCreateInstance] failed: instanceDriver.Start instanceID=%s err=%v", inst.ID(), err)
		return err
	}
	log.Printf("[ExecuteCreateInstance] instance started instanceID=%s", inst.ID())

	log.Printf("[ExecuteCreateInstance] loading ingresses instanceID=%s", inst.ID())
	ingresses, err := i.ingressRepo.FindByInstanceID(ctx, inst.ID())
	if err != nil {
		log.Printf("[ExecuteCreateInstance] failed: FindByInstanceID for ingresses instanceID=%s err=%v", inst.ID(), err)
		return err
	}
	log.Printf("[ExecuteCreateInstance] ingresses loaded count=%d instanceID=%s", len(ingresses), inst.ID())

	log.Printf("[ExecuteCreateInstance] applying routes to gateway instanceID=%s count=%d", inst.ID(), len(ingresses))
	if err := i.gatewayDriver.ApplyRoutes(ctx, ingresses); err != nil {
		log.Printf("[ExecuteCreateInstance] failed: gatewayDriver.ApplyRoutes instanceID=%s err=%v", inst.ID(), err)
		inst.MarkAsError(compute.ErrInStarting)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}
	log.Printf("[ExecuteCreateInstance] routes applied to gateway instanceID=%s", inst.ID())

	// 最終的な状態をDBに保存する
	log.Printf("[ExecuteCreateInstance] marking running instanceID=%s", inst.ID())
	err = i.uow.Do(ctx, func(ctx context.Context) error {
		var uowErr error
		// 最新状態をDBから取得する
		inst, uowErr = i.instanceRepo.FindByID(ctx, payload.InstanceID)
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: FindByID (for running) instanceID=%s err=%v", payload.InstanceID, uowErr)
			return compute.ErrInstanceNotFound
		}
		uowErr = inst.MarkAsRunning() // 状態を「起動中」に遷移させる
		if uowErr != nil {
			log.Printf("[ExecuteCreateInstance] failed: MarkAsRunning instanceID=%s err=%v", inst.ID(), uowErr)
			return compute.ErrInvalidInstanceStatus
		}
		log.Printf("[ExecuteCreateInstance] marked running instanceID=%s", inst.ID())
		return i.instanceRepo.Save(ctx, inst)
	})
	if err != nil {
		log.Printf("[ExecuteCreateInstance] failed: uow (running) err=%v", err)
		inst.MarkAsError(compute.ErrInStarting)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}
	log.Printf("[ExecuteCreateInstance] completed instanceID=%s status=%s", inst.ID(), inst.Status())

	return nil
}
