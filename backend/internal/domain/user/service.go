package user

import "context"

type UserIdentity struct {
	ID          UserID
	DisplayName string
	Permissions []Permission
}

type IdentityService interface {
	GetIdentity(ctx context.Context, sub string) (*UserIdentity, error)
}
