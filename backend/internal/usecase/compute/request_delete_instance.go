package compute

import (
	"context"
	"encoding/json"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/user"
	"example.com/m/internal/usecase"
)

type RequestDeleteInstanceInput struct {
	InstanceID compute.InstanceID
}

type RequestDeleteInstanceUseCase interface {
	Execute(ctx context.Context, input RequestDeleteInstanceInput) error
}

type requestDeleteInstanceInteractor struct {
	instanceRepo compute.InstanceRepository
	publisher    usecase.JobPublisher
}

func NewRequestDeleteInstanceInteractor() RequestDeleteInstanceUseCase {
	return &requestDeleteInstanceInteractor{}
}

func (i *requestDeleteInstanceInteractor) Execute(ctx context.Context, input RequestDeleteInstanceInput) error {
	usr, ok := user.FromContext(ctx)
	if !ok {
		return user.ErrUserNotInContext
	}

	inst, err := i.instanceRepo.FindByID(ctx, input.InstanceID)
	if err != nil {
		return err
	}

	if inst.OwnerID() != usr.ID() {
		if !usr.HasPermission(user.PermissionInstanceDeleteAll) {
			return user.ErrNoPermission
		}
	} else {
		if !usr.HasPermission(user.PermissionInstanceDelete) {
			return user.ErrNoPermission
		}
	}

	if err := inst.MarkAsDeleting(); err != nil {
		return err
	}

	if err := i.instanceRepo.Save(ctx, inst); err != nil {
		return err
	}

	payload := DeleteInstancePayload{
		InstanceID: inst.ID(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if err := i.publisher.Publish(ctx, usecase.JobTypeDeleteInstance, payloadBytes); err != nil {
		inst.MarkAsError(compute.ErrInDeleting)
		_ = i.instanceRepo.Save(ctx, inst)
		return err
	}

	return nil
}
