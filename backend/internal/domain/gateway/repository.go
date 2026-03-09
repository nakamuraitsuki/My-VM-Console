package gateway

import "context"

type Repository interface {
	FindByID(ctx context.Context, id IngressID) (*IngressRoute, error)
	FindByOwnerID(ctx context.Context, ownerID string) ([]*IngressRoute, error)
	Save(ctx context.Context, route *IngressRoute) error
	Delete(ctx context.Context, id IngressID) error
}
