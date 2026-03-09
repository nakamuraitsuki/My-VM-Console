package gateway

import "context"

type IngressDriver interface {
	ApplyRoute(ctx context.Context, route *IngressRoute) error
	RemoveRoute(ctx context.Context, route *IngressRoute) error
}
