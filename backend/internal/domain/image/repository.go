package image

import (
	"context"
	"errors"
)

var (
	ErrImageNotFound = errors.New("image not found")
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Image, error)
	FindByAlias(ctx context.Context, alias string) (*Image, error)
	FindByID(ctx context.Context, id ImageID) (*Image, error)
	Save(ctx context.Context, img *Image) error
}
