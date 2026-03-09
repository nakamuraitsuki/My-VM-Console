package compute

import "context"

type ComputeDriver interface {
	Create(ctx context.Context, instance *Instance) error
	Start(ctx context.Context, id InstanceID) error
	Stop(ctx context.Context, id InstanceID) error
	Terminate(ctx context.Context, id InstanceID) error
	GetRealStatus(ctx context.Context, id InstanceID) (InstanceStatus, error)
}
