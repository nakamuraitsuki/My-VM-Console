package compute

import (
	"context"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/gateway"
	"example.com/m/internal/domain/storage"
	"example.com/m/internal/domain/user"
)

type GetInstanceInput struct {
	InstanceID compute.InstanceID
}

type GetInstanceOutput struct {
	Instance  *compute.Instance
	Ingresses []*gateway.IngressRoute
	Volumes   *storage.Volume
}

type GetInstanceUseCase interface {
	Execute(ctx context.Context, input GetInstanceInput) (*GetInstanceOutput, error)
}

type getInstanceInteractor struct {
	instanceRepository compute.InstanceRepository
	gatewayRepository  gateway.Repository
	storageRepository  storage.Repository
}

func NewGetInstanceInteractor(
	instanceRepository compute.InstanceRepository,
	gatewayRepository gateway.Repository,
	storageRepository storage.Repository,
) GetInstanceUseCase {
	return &getInstanceInteractor{
		instanceRepository: instanceRepository,
		gatewayRepository:  gatewayRepository,
		storageRepository:  storageRepository,
	}
}

func (i *getInstanceInteractor) Execute(ctx context.Context, input GetInstanceInput) (*GetInstanceOutput, error) {
	usr, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotInContext
	}

	inst, err := i.instanceRepository.FindByID(ctx, input.InstanceID)
	if err != nil {
		return nil, err
	}

	if inst.OwnerID() != usr.ID() {
		usr.HasPermission(user.PermissionInstanceRead)
	}

	ingresses, err := i.gatewayRepository.FindByInstanceID(ctx, inst.ID())
	if err != nil {
		return nil, err
	}

	volumes, err := i.storageRepository.FindByID(ctx, inst.RootVolumeID())
	if err != nil {
		return nil, err
	}

	return &GetInstanceOutput{
		Instance:  inst,
		Ingresses: ingresses,
		Volumes:   volumes,
	}, nil
}
