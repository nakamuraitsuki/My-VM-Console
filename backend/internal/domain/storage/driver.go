package storage

import "context"

type StorageDriver interface {
	CreateVolume(ctx context.Context, vol *Volume) error
	DeleteVolume(ctx context.Context, vol *Volume) error
}
