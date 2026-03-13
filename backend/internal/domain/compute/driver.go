package compute

import (
	"context"

	"example.com/m/internal/domain/image"
)

type ComputeDriver interface {
	Create(ctx context.Context, instance *Instance, image *image.Image) error
	Start(ctx context.Context, id InstanceID) error
	Stop(ctx context.Context, id InstanceID) error
	Terminate(ctx context.Context, id InstanceID) error
	GetRealStatus(ctx context.Context, id InstanceID) (InstanceStatus, error)
}
