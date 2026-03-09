package compute

import "context"

type InstanceRepository interface {
	FindByID(ctx context.Context, id InstanceID) (*Instance, error)
	FindByOwnerID(ctx context.Context, ownerID string) ([]*Instance, error)
	Save(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context, id InstanceID) error
}
