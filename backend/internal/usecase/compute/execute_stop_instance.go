package compute

import (
	"context"

	"example.com/m/internal/domain/compute"
)

type StopInstancePayload struct {
	InstanceID compute.InstanceID
}

type ExecuteStopInstanceUseCase interface {
	Execute(ctx context.Context, payload StopInstancePayload) error
}

type executeStopInstanceInteractor struct {
	instanceRepo compute.InstanceRepository
	driver 	 compute.ComputeDriver
}

func NewExecuteStopInstanceInteractor() ExecuteStopInstanceUseCase {
	return &executeStopInstanceInteractor{}
}

func (i *executeStopInstanceInteractor) Execute(ctx context.Context, payload StopInstancePayload) error {
	inst, err := i.instanceRepo.FindByID(ctx, payload.InstanceID)
	if err != nil {
		return err
	}

	if err := inst.MarkAsStopped(); err != nil {
		return err
	}
	
	if err := i.instanceRepo.Save(ctx, inst); err != nil {
		return err
	}

	if err := i.driver.Stop(ctx, payload.InstanceID); err != nil {
		inst.MarkAsError(compute.ErrInStopping)
		_ = i.instanceRepo.Save(ctx, inst)
		return err
	}

	return nil
}