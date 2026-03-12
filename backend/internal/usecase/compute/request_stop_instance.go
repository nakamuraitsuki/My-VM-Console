package compute

import (
	"context"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/user"
	"example.com/m/internal/usecase"
)

type RequestStopInstanceInput struct {
	InstanceID compute.InstanceID
}

type StopInstancePayload struct {
	InstanceID compute.InstanceID
}

type RequestStopInstanceUseCase interface {
	Execute(ctx context.Context, input RequestStopInstanceInput) error
}

type requestStopInstanceInteractor struct {
	instanceRepo compute.InstanceRepository
	publisher    usecase.JobPublisher
	uow          usecase.UnitOfWork
}

func NewRequestStopInstanceInteractor() RequestStopInstanceUseCase {
	return &requestStopInstanceInteractor{}
}

func (i *requestStopInstanceInteractor) Execute(ctx context.Context, input RequestStopInstanceInput) error {
	usr, ok := user.FromContext(ctx)
	if !ok {
		return user.ErrUserNotInContext
	}

	inst, err := i.instanceRepo.FindByID(ctx, input.InstanceID)
	if err != nil {
		return err
	}

	if inst.OwnerID() != usr.ID() {
		if !usr.HasPermission(user.PermissionInstanceStopAll) {
			return user.ErrNoPermission
		}
	} else {
		if !usr.HasPermission(user.PermissionInstanceStop) {
			return user.ErrNoPermission
		}
	}

	return i.uow.Do(ctx, func(ctx context.Context) error {
		if err := inst.MarkAsStopping(); err != nil {
			return err
		}

		if err := i.instanceRepo.Save(ctx, inst); err != nil {
			return err
		}

		payload := StopInstancePayload{
			InstanceID: inst.ID(),
		}

		if err := i.publisher.Publish(ctx, usecase.JobTypeStopInstance, payload); err != nil {
			return err
		}
		
		return nil
	})
}
