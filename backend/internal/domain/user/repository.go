package user

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserPersistentData struct {
	ID    UserID
	Quota UsageQuota
	Status UserStatus
	ErrorPhase *FailedPhase // エラー理由（エラー状態のときのみ値が入る）
}

type UserRepository interface {
	FindByID(ctx context.Context, id UserID) (*UserPersistentData, error)
	Save(ctx context.Context, user *UserPersistentData) error
	Delete(ctx context.Context, id UserID) error
}
