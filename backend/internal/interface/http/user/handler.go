package user

import (
	"example.com/m/internal/infrastructure/auth/oidc"
	"example.com/m/internal/usecase/user"
)

type Handler struct {
	oidcConfig        *oidc.OIDCConfig
	idTokenVerifier   oidc.IDTokenVerifier
	ensureUserUseCase user.EnsureUserUseCase
}

func NewHandler(
	oidcConfig *oidc.OIDCConfig,
	idTokenVerifier oidc.IDTokenVerifier,
	ensureUserUseCase user.EnsureUserUseCase,
) *Handler {
	return &Handler{
		oidcConfig:        oidcConfig,
		idTokenVerifier:   idTokenVerifier,
		ensureUserUseCase: ensureUserUseCase,
	}
}
