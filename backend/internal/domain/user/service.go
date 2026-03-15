package user

import "context"

type UserIdentity struct {
	ID              UserID
	DisplayName     string
	ProfileImageURL string
	Permissions     []Permission
}

type IdentityService interface {
	GetIdentity(ctx context.Context, token string) (*UserIdentity, error)
}
