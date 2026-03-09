package compute

import (
	"errors"

	"example.com/m/internal/domain/image"
	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/storage"
	"example.com/m/internal/domain/user"
)

var (
	ErrInstanceAlreadyRunning = errors.New("instance is already running")
	ErrInstanceNotRunnning    = errors.New("instance is not running")
)

type InstanceID string

type InstanceStatus string

const (
	StatusStopped  InstanceStatus = "stopped"
	StatusRunnning InstanceStatus = "running"
	StatePending   InstanceStatus = "pending"
	StatusDeleting InstanceStatus = "deleting"
	StatusError    InstanceStatus = "error"
)

type Instance struct {
	id      InstanceID
	name    string
	ownerID user.UserID
	status  InstanceStatus

	// スペック
	cpu      int
	memoryMB int

	// Image
	imageID image.ImageID

	// Network
	subenetID network.SubnetID
	privateIP string

	// Storage
	rootVolumeID storage.VolumeID
}

// --- Constructor ---
func NewInstance(
	id InstanceID,
	name string,
	owner user.UserID,
	cpu, memoryMB int,
	imageID image.ImageID,
	subnetID network.SubnetID,
	privateIP string,
	rootVolumeID storage.VolumeID,
) *Instance {
	return &Instance{
		id:           id,
		name:         name,
		ownerID:      owner,
		status:       StatusStopped,
		cpu:          cpu,
		memoryMB:     memoryMB,
		imageID:      imageID,
		subenetID:    subnetID,
		privateIP:    privateIP,
		rootVolumeID: rootVolumeID,
	}
}

// --- Getters ---
func (i *Instance) ID() InstanceID                 { return i.id }
func (i *Instance) Name() string                   { return i.name }
func (i *Instance) OwnerID() user.UserID           { return i.ownerID }
func (i *Instance) Status() InstanceStatus         { return i.status }
func (i *Instance) ImageID() image.ImageID         { return i.imageID }
func (i *Instance) SubnetID() network.SubnetID     { return i.subenetID }
func (i *Instance) PrivateIP() string              { return i.privateIP }
func (i *Instance) RootVolumeID() storage.VolumeID { return i.rootVolumeID }
