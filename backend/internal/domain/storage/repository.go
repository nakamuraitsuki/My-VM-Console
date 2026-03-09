package storage

import "context"

type Repository interface {
	FindByID(ctx context.Context, id VolumeID) (*Volume, error)
	Save(ctx context.Context, volume *Volume) error
	Delete(ctx context.Context, id VolumeID) error
}
