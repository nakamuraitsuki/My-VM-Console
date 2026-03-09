package storage

type VolumeID string

type Volume struct {
	id     VolumeID
	name   string
	sizeGB int
	pool   string // zfs, btrfs, dir など
	owner  string
}
