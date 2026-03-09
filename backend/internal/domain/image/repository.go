package image

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*Image, error)
	FindByAlias(ctx context.Context, alias string) (*Image, error)
	Save(ctx context.Context, img *Image) error
}
