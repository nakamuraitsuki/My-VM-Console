package user

import "errors"

// --- Errors ---
var (
	ErrQuotaExceeded = errors.New("resource quota exceeded")
	ErrNoPermission  = errors.New("permission denied")
	ErrInvalidStatus  = errors.New("invalid user status for this operation")
)

type UserID string
type Permission string

// --- Permissions ---
const (
	PermissionAll            Permission = "*"
	PermissionInstanceRead   Permission = "instance:read"
	PermissionInstanceCreate Permission = "instance:create"
	PermissionInstanceUpdate Permission = "instance:update"
	PermissionInstanceDelete Permission = "instance:delete"
	PermissionNetworkManage  Permission = "network:manage"
)

type UserStatus string

// --- User Status ---
const (
	UserStatusPending      UserStatus = "pending"
	UserStatusInitializing UserStatus = "initializing"
	UserStatusActive       UserStatus = "active"
	UserStatusFailed       UserStatus = "failed"
)

type FailedPhase string

// --- Failed Phase ---
const (
	FailedInPending      FailedPhase = "failed in pending"
	FailedInInitializing FailedPhase = "failed in initializing"
)

type UsageQuota struct {
	MaxInstance int
	MaxCPU      int
	MaxMemory   int
}

type User struct {
	id          UserID
	displayName string
	permissions []Permission
	quota       UsageQuota
	status      UserStatus
	errorPhase *FailedPhase // エラー理由（エラー状態のときのみ値が入る）
}

// --- Constructor ---
func NewUser(id UserID, name string, perms []Permission, quota UsageQuota, status UserStatus, errorPhase *FailedPhase) *User {
	return &User{
		id:          id,
		displayName: name,
		permissions: perms,
		quota:       quota,
		status:      status,
		errorPhase: errorPhase,
	}
}

// --- Getters ---
func (u *User) ID() UserID          { return u.id }
func (u *User) DisplayName() string { return u.displayName }
func (u *User) Quota() UsageQuota   { return u.quota }

// --- Setters / Domain Logic ---
func (u *User) UpdateQuota(newQuota UsageQuota) {
	u.quota = newQuota
}

// --- Business Logic ---

// 引数は必要とされる権限。
func (u *User) HasPermission(perm Permission) bool {
	for _, p := range u.permissions {
		if p == PermissionAll || p == perm {
			return true
		}
	}
	return false
}

func (u *User) CanAllocateInstance(currentInstances, requestedCPU int) bool {
	if currentInstances+1 > u.quota.MaxInstance {
		return false
	}
	if requestedCPU > u.quota.MaxCPU {
		return false
	}
	return true
}

// 冪等性の確保のため、pending, failed && errorPhase=initializingのときのみネットワーク作成処理を実行する
func (u *User) MarkAsInitializing() error {
	if u.status == UserStatusActive {
		return ErrInvalidStatus
	}
	if u.status == UserStatusFailed && (u.errorPhase == nil || *u.errorPhase != FailedInPending) {
		return ErrInvalidStatus
	}
	u.status = UserStatusInitializing
	u.errorPhase = nil // エラー理由はクリア
	return nil
}

func (u *User) MarkAsActive() error {
	if u.status != UserStatusInitializing {
		return ErrInvalidStatus
	}
	u.status = UserStatusActive
	return nil
}

func (u *User) MarkAsFailed(phase FailedPhase) UserPersistentData {
	u.status = UserStatusFailed
	u.errorPhase = &phase
	return UserPersistentData{
		ID: u.id,
		Quota: u.quota,
		Status: u.status,
		ErrorPhase: u.errorPhase,
	}
}