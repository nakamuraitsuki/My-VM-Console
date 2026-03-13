package user

import (
	"context"
	"errors"
)

var (
	ErrUserNotInContext = errors.New("user not found in context")
)

type ctxKey struct{}

func WithContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, ctxKey{}, user)
}

func FromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(ctxKey{}).(*User)
	return user, ok
}
