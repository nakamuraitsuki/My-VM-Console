package compute

import (
	"context"

	"example.com/m/internal/domain/image"
)

type ComputeDriver interface {
	Create(ctx context.Context, instance *Instance, image *image.Image) error
	Start(ctx context.Context, instance *Instance) error
	Stop(ctx context.Context, instance *Instance) error
	Terminate(ctx context.Context, instance *Instance) error
	GetRealStatus(ctx context.Context, instance *Instance) (InstanceStatus, error)

	AuthorizePublicKey(ctx context.Context, instance *Instance, publicKey string) error
	RevokePublicKey(ctx context.Context, instance *Instance, publicKey string) error
}
