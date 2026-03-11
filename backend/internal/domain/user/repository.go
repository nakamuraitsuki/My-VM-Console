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
}

type UserRepository interface {
	FindByID(ctx context.Context, id UserID) (*UserPersistentData, error)
	Save(ctx context.Context, user *UserPersistentData) error
	Delete(ctx context.Context, id UserID) error
}
