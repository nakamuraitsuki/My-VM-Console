package user

import "errors"

// --- Errors ---
var (
	ErrQuotaExceeded = errors.New("resource quota exceeded")
	ErrNoPermission  = errors.New("permission denied")
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
}

// --- Constructor ---
func NewUser(id UserID, name string, perms []Permission, quota UsageQuota) *User {
	return &User{
		id:          id,
		displayName: name,
		permissions: perms,
		quota:       quota,
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

func (u *User) CanAllocateResources(currentInstances, requestedCPU, requestedMemory int) error {
	if currentInstances+1 > u.quota.MaxInstance {
		return ErrQuotaExceeded
	}
	if requestedCPU > u.quota.MaxCPU {
		return ErrQuotaExceeded
	}
	if requestedMemory > u.quota.MaxMemory {
		return ErrQuotaExceeded
	}
	return nil
}

func (u *User) CanAllocateInstance(currentInstances, requestedCPU int) error {
	if currentInstances+1 > u.quota.MaxInstance {
		return ErrQuotaExceeded
	}
	if requestedCPU > u.quota.MaxCPU {
		return ErrQuotaExceeded
	}
	return nil
}
