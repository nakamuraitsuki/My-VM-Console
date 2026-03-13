package storage

import "github.com/google/uuid"

type VolumeID string

type Volume struct {
	id     VolumeID
	name   string
	sizeGB int
	pool   string // zfs, btrfs, dir など
	owner  string
}

func NewVolume(id VolumeID, name string, sizeGB int, pool string, owner string) *Volume {
	return &Volume{
		id:     id,
		name:   name,
		sizeGB: sizeGB,
		pool:   pool,
		owner:  owner,
	}
}

func NewID() VolumeID {
	return VolumeID("vol-" + uuid.New().String())
}

func (v *Volume) ID() VolumeID  { return v.id }
func (v *Volume) Name() string  { return v.name }
func (v *Volume) SizeGB() int   { return v.sizeGB }
func (v *Volume) Pool() string  { return v.pool }
func (v *Volume) Owner() string { return v.owner }
