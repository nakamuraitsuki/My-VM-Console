package gateway

import "context"

type IngressDriver interface {
	ApplyRoutes(ctx context.Context, routes []*IngressRoute) error
	RemoveRoute(ctx context.Context, route *IngressRoute) error
}
