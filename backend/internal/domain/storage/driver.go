package storage

import (
	"context"

	"example.com/m/internal/domain/network"
)

type StorageDriver interface {
	CreateVolume(ctx context.Context, vpcID network.VPCID, vol *Volume) error
	DeleteVolume(ctx context.Context, vpcID network.VPCID, vol *Volume) error
}
