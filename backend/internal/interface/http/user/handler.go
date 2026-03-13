package user

import (
	"example.com/m/internal/infrastructure/auth/oidc"
	"example.com/m/internal/usecase/user"
)

type Handler struct {
	oidcConfig            *oidc.OIDCConfig
	idTokenVerifier       oidc.IDTokenVerifier
	ensureUserUseCase     user.EnsureUserUseCase
	listMyInstanceUseCase user.ListMyInstanceUseCase
}

func NewHandler(
	oidcConfig *oidc.OIDCConfig,
	idTokenVerifier oidc.IDTokenVerifier,
	ensureUserUseCase user.EnsureUserUseCase,
	listMyInstanceUseCase user.ListMyInstanceUseCase,
) *Handler {
	return &Handler{
		oidcConfig:            oidcConfig,
		idTokenVerifier:       idTokenVerifier,
		ensureUserUseCase:     ensureUserUseCase,
		listMyInstanceUseCase: listMyInstanceUseCase,
	}
}
