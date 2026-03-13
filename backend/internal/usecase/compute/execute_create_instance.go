package compute

import (
	"context"
	"errors"
	"fmt"

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

	if err := i.uow.Do(ctx, func(ctx context.Context) error {
		var uowErr error
		inst, uowErr = i.instanceRepo.FindByID(ctx, payload.InstanceID)
		if uowErr != nil {
			return compute.ErrInstanceNotFound
		}

		img, uowErr = i.imgRepo.FindByID(ctx, inst.ImageID())
		if uowErr != nil {
			return uowErr
		}
		if img == nil {
			return image.ErrImageNotFound
		}

		vpc, uowErr = i.networkRepo.FindVPCByUserID(ctx, inst.OwnerID())
		if uowErr != nil {
			return uowErr
		}
		if vpc == nil {
			return fmt.Errorf("VPC not found for user: %s", inst.OwnerID())
		}

		uowErr = inst.MarkAsCreating() // 状態を「作成中」に遷移させる
		if uowErr != nil {
			return compute.ErrInvalidInstanceStatus
		}

		return i.instanceRepo.Save(ctx, inst)
	}); err != nil {
		if errors.Is(err, compute.ErrInstanceNotFound) {
			return compute.ErrInstanceNotFound
		}
		inst.MarkAsError(compute.ErrInPending)
		return err
	}

	//　storage
	volumeID := inst.RootVolumeID()
	volume, err := i.storageRepo.FindByID(ctx, volumeID)
	if err != nil {
		return err
	}
	err = i.storageDriver.CreateVolume(ctx, vpc.ID(), volume)
	if err != nil {
		inst.MarkAsError(compute.ErrInCreating)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}

	// instance
	err = i.instanceDriver.Create(ctx, inst, img)
	if err != nil {
		inst.MarkAsError(compute.ErrInCreating)
		_ = i.instanceRepo.Save(ctx, inst)
		_ = i.storageDriver.DeleteVolume(ctx, vpc.ID(), volume) // 作成したVolumeを削除
		return err
	}

	// 遷移
	err = i.uow.Do(ctx, func(ctx context.Context) error {
		var uowErr error
		// 最新状態をDBから取得する
		inst, uowErr = i.instanceRepo.FindByID(ctx, payload.InstanceID)
		if uowErr != nil {
			return compute.ErrInstanceNotFound
		}
		uowErr = inst.MarkAsStarting() // 状態を「起動中」に遷移させる
		if uowErr != nil {
			return compute.ErrInvalidInstanceStatus
		}
		return i.instanceRepo.Save(ctx, inst)
	})
	if err != nil {
		inst.MarkAsError(compute.ErrInCreating)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}

	err = i.instanceDriver.Start(ctx, inst.ID())
	if err != nil {
		return err
	}

	ingresses, err := i.ingressRepo.FindByInstanceID(ctx, inst.ID())
	if err != nil {
		return err
	}

	if err := i.gatewayDriver.ApplyRoutes(ctx, ingresses); err != nil {
		inst.MarkAsError(compute.ErrInStarting)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}

	// 最終的な状態をDBに保存する
	err = i.uow.Do(ctx, func(ctx context.Context) error {
		var uowErr error
		// 最新状態をDBから取得する
		inst, uowErr = i.instanceRepo.FindByID(ctx, payload.InstanceID)
		if uowErr != nil {
			return compute.ErrInstanceNotFound
		}
		uowErr = inst.MarkAsRunning() // 状態を「起動中」に遷移させる
		if uowErr != nil {
			return compute.ErrInvalidInstanceStatus
		}
		return i.instanceRepo.Save(ctx, inst)
	})
	if err != nil {
		inst.MarkAsError(compute.ErrInStarting)
		_ = i.instanceRepo.Save(ctx, inst) // エラー状態を保存
		return err
	}

	return nil
}
