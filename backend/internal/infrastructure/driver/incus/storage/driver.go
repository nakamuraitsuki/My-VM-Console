package storage

import (
	"context"
	"fmt"

	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/storage"
	incus "github.com/lxc/incus/v6/client"
	"github.com/lxc/incus/v6/shared/api"
)

type driver struct {
	client incus.InstanceServer
}

func NewDriver(c incus.InstanceServer) storage.StorageDriver {
	return &driver{
		client: c,
	}
}

func (d *driver) CreateVolume(ctx context.Context, vpcID network.VPCID, vol *storage.Volume) error {
	pc := d.client.UseProject(string(vpcID))

	put := api.StorageVolumePut{
		Config: map[string]string{
			"size": fmt.Sprintf("%dGiB", vol.SizeGB()),
		},
		Description: fmt.Sprintf("Volume for Instance %s", vol.Owner()),
	}

	req := api.StorageVolumesPost{
		StorageVolumePut: put,
		Name:             string(vol.ID()),
		Type:             "custom",
	}

	poolName := "default" // Incusのストレージプール名。将来的に変更できるようにするかも
	err := pc.CreateStoragePoolVolume(poolName, req)
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) DeleteVolume(ctx context.Context, vpcID network.VPCID, vol *storage.Volume) error {
	pc := d.client.UseProject(string(vpcID))

	poolName := "default"

	err := pc.DeleteStoragePoolVolume(poolName, "custom", string(vol.ID()))
	if err != nil {
		return err
	}

	return nil
}
