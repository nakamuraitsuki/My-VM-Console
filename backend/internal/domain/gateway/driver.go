package gateway

import "context"

type IngressDriver interface {
	ApplyRoutes(ctx context.Context, routes []*IngressRoute) error
	RemoveRoutes(ctx context.Context, routes []*IngressRoute) error
}
