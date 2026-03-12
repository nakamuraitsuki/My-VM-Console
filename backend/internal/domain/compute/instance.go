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
	ErrInstanceNotRunning     = errors.New("instance is not running")
	ErrInvalidInstanceStatus  = errors.New("invalid instance status")
)

type InstanceID string

type InstanceStatus string

const (
	StatusPending  InstanceStatus = "pending"
	StatusCreating InstanceStatus = "creating"
	StatusStarting InstanceStatus = "starting"
	StatusRunning  InstanceStatus = "running"
	StatusStopping InstanceStatus = "stopping"
	StatusStopped  InstanceStatus = "stopped"
	StatusDeleting InstanceStatus = "deleting"
	StatusError    InstanceStatus = "error"
)

type ErrPhase string

const (
	ErrInPending  ErrPhase = "error in pending"
	ErrInCreating ErrPhase = "error in creating"
	ErrInStarting ErrPhase = "error in starting"
	ErrInStopping ErrPhase = "error in stopping"
	ErrInDeleting ErrPhase = "error in deleting"
)

type Instance struct {
	id         InstanceID
	name       string
	ownerID    user.UserID
	status     InstanceStatus
	errorPhase *ErrPhase // エラー理由（エラー状態のときのみ値が入る）

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
	status InstanceStatus,
	errorPhase *ErrPhase,
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
		status:       status,
		errorPhase:   errorPhase,
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
func (i *Instance) CPU() int                       { return i.cpu }
func (i *Instance) MemoryMB() int                  { return i.memoryMB }
func (i *Instance) ErrPhase() *ErrPhase            { return i.errorPhase }

// --- Setters ---

// 状態確認ののち、「作成中」状態に遷移させる
func (i *Instance) MarkAsCreating() error {
	if i.status == StatusPending {
		i.status = StatusCreating
		return nil
	}
	if i.status == StatusError && i.errorPhase != nil && *i.errorPhase == ErrInPending {
		i.status = StatusCreating
		i.errorPhase = nil
		return nil
	}
	return ErrInvalidInstanceStatus
}

func (i *Instance) MarkAsStarting() error {
	if i.status == StatusCreating || i.status == StatusStopped {
		i.status = StatusStarting
		return nil
	}
	if i.status == StatusError && i.errorPhase != nil && (*i.errorPhase == ErrInCreating || *i.errorPhase == ErrInStarting) {
		i.status = StatusStarting
		i.errorPhase = nil
		return nil
	}
	return ErrInvalidInstanceStatus
}

func (i *Instance) MarkAsRunning() error {
	if i.status == StatusStarting {
		i.status = StatusRunning
		return nil
	}
	if i.status == StatusError && i.errorPhase != nil && *i.errorPhase == ErrInStarting {
		i.status = StatusRunning
		i.errorPhase = nil
		return nil
	}
	return ErrInvalidInstanceStatus
}

func (i *Instance) MarkAsStopping() error {
	// 以上終了を許可するためのStatusError許可
	if i.status == StatusRunning || i.status == StatusError {
		i.status = StatusStopping
		return nil
	}
	return ErrInvalidInstanceStatus
}

func (i *Instance) MarkAsStopped() error {
	if i.status == StatusStopping {
		i.status = StatusStopped
		return nil
	}
	if i.status == StatusError && i.errorPhase != nil && *i.errorPhase == ErrInStopping {
		i.status = StatusStopped
		i.errorPhase = nil
		return nil
	}
	return ErrInvalidInstanceStatus
}

func (i *Instance) MarkAsDeleting() error {
	if i.status == StatusStopped || i.status == StatusError {
		i.status = StatusDeleting
		return nil
	}
	return ErrInvalidInstanceStatus
}

// エラー状態に遷移させる
func (i *Instance) MarkAsError(errorPhase ErrPhase) {
	i.status = StatusError
	i.errorPhase = &errorPhase
}
