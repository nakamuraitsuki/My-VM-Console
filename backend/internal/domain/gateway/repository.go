package gateway

import (
	"context"

	"example.com/m/internal/domain/compute"
)

type Repository interface {
	FindByID(ctx context.Context, id IngressID) (*IngressRoute, error)
	FindByInstanceID(ctx context.Context, instanceID compute.InstanceID) ([]*IngressRoute, error)
	FindByOwnerID(ctx context.Context, ownerID string) ([]*IngressRoute, error)
	Save(ctx context.Context, route *IngressRoute) error
	Delete(ctx context.Context, id IngressID) error
}
