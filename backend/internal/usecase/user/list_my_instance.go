package user

import (
	"context"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/user"
)

type ListMyInstanceUseCase interface {
	Execute(ctx context.Context) ([]*compute.Instance, error)
}

type ListMyInstanceInteractor struct {
	InstanceRepository compute.InstanceRepository
}

func NewListMyInstanceInteractor(instanceRepository compute.InstanceRepository) *ListMyInstanceInteractor {
	return &ListMyInstanceInteractor{
		InstanceRepository: instanceRepository,
	}
}

func (i *ListMyInstanceInteractor) Execute(ctx context.Context) ([]*compute.Instance, error) {
	usr, ok := user.FromContext(ctx)
	if !ok {
		return nil, user.ErrUserNotInContext
	}

	instances, err := i.InstanceRepository.FindByOwnerID(ctx, usr.ID())
	if err != nil {
		return nil, err
	}
	return instances, nil
}
