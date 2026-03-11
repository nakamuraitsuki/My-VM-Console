package compute

import (
	"context"
	"errors"

	"example.com/m/internal/domain/user"
)

var (
	ErrInstanceNotFound = errors.New("instance not found")
)

type InstanceRepository interface {
	FindByID(ctx context.Context, id InstanceID) (*Instance, error)
	FindByOwnerID(ctx context.Context, ownerID user.UserID) ([]*Instance, error)
	Save(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context, id InstanceID) error
}
