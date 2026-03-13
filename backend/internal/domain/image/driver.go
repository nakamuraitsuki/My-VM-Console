package image

import "context"

type ImageDriver interface {
	PullImage(ctx context.Context, alias string) error
	Exists(ctx context.Context, alias string) (bool, error)
}
