package user

import "context"

type ctxKey struct{}

func WithContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, ctxKey{}, user)
}

func FromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(ctxKey{}).(*User)
	return user, ok
}