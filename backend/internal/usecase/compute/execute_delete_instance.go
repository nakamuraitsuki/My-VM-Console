package compute

import (
	"context"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/gateway"
	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/storage"
	"example.com/m/internal/usecase"
)

type DeleteInstancePayload struct {
	InstanceID compute.InstanceID
}

type ExecuteDeleteInstanceUseCase interface {
	Execute(ctx context.Context, payload DeleteInstancePayload) error
}

type executeDeleteInstanceInteractor struct {
	instanceRepo   compute.InstanceRepository
	networkRepo    network.Repository
	volumeRepo     storage.Repository
	gatewayRepo    gateway.Repository
	instanceDriver compute.ComputeDriver
	storageDriver  storage.StorageDriver
	networkDriver  network.NetworkDriver
	gatewayDriver  gateway.IngressDriver
	uow            usecase.UnitOfWork
}

func NewExecuteDeleteInstanceInteractor(
	instanceRepo compute.InstanceRepository,
	networkRepo network.Repository,
	volumeRepo storage.Repository,
	gatewayRepo gateway.Repository,
	instanceDriver compute.ComputeDriver,
	storageDriver storage.StorageDriver,
	networkDriver network.NetworkDriver,
	gatewayDriver gateway.IngressDriver,
	uow usecase.UnitOfWork,
) ExecuteDeleteInstanceUseCase {
	return &executeDeleteInstanceInteractor{
		instanceRepo:   instanceRepo,
		networkRepo:    networkRepo,
		volumeRepo:     volumeRepo,
		gatewayRepo:    gatewayRepo,
		instanceDriver: instanceDriver,
		storageDriver:  storageDriver,
		networkDriver:  networkDriver,
		gatewayDriver:  gatewayDriver,
		uow:            uow,
	}
}

// NOTE: Driverの実装においてべき等性が確保されている前提で組まれている。
func (i *executeDeleteInstanceInteractor) Execute(ctx context.Context, payload DeleteInstancePayload) error {
	inst, err := i.instanceRepo.FindByID(ctx, payload.InstanceID)
	if err != nil {
		return err
	}

	// 各種リソースの物理的削除
	// Root Volumeの削除
	vol, err := i.volumeRepo.FindByID(ctx, inst.RootVolumeID())
	if err != nil {
		return err
	}
	if err := i.storageDriver.DeleteVolume(ctx, vol); err != nil {
		inst.MarkAsError(compute.ErrInDeleting)
		_ = i.instanceRepo.Save(ctx, inst)
		return err
	}

	// gatewayのルート削除
	ingresses, err := i.gatewayRepo.FindByInstanceID(ctx, inst.ID())
	if err != nil {
		inst.MarkAsError(compute.ErrInDeleting)
		_ = i.instanceRepo.Save(ctx, inst)
		return err
	}
	if len(ingresses) > 0 {
		if err := i.gatewayDriver.RemoveRoutes(ctx, ingresses); err != nil {
			inst.MarkAsError(compute.ErrInDeleting)
			_ = i.instanceRepo.Save(ctx, inst)
			return err
		}
	}

	if err := i.instanceDriver.Terminate(ctx, inst.ID()); err != nil {
		inst.MarkAsError(compute.ErrInDeleting)
		_ = i.instanceRepo.Save(ctx, inst)
		return err
	}

	// 論理削除
	uowErr := i.uow.Do(ctx, func(ctx context.Context) error {
		if err := i.instanceRepo.Delete(ctx, inst.ID()); err != nil {
			return err
		}
		if err := i.volumeRepo.Delete(ctx, inst.RootVolumeID()); err != nil {
			return err
		}
		if err := i.networkRepo.DeleteLease(ctx, string(inst.ID())); err != nil {
			return err
		}
		ids := make([]gateway.IngressID, len(ingresses))
		for i, ing := range ingresses {
			ids[i] = ing.ID()
		}
		if err := i.gatewayRepo.DeleteBulk(ctx, ids); err != nil {
			return err
		}
		return nil
	})
	if uowErr != nil {
		inst.MarkAsError(compute.ErrInDeleting)
		_ = i.instanceRepo.Save(ctx, inst)
		return uowErr
	}
	return nil
}
