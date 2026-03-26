package compute

import (
	"example.com/m/internal/usecase/compute"
	"example.com/m/internal/usecase/user"
)

type Handler struct {
	reqCreateUseCase    compute.RequestCreateInstanceUseCase
	ensureUserUseCase   user.EnsureUserUseCase
	authorizeKeyUseCase compute.AuthorizePubkeyUseCase
}

func NewHandler(
	reqCreateUseCase compute.RequestCreateInstanceUseCase,
	ensureUserUseCase user.EnsureUserUseCase,
	authorizeKeyUseCase compute.AuthorizePubkeyUseCase,
) *Handler {
	return &Handler{
		reqCreateUseCase:    reqCreateUseCase,
		ensureUserUseCase:   ensureUserUseCase,
		authorizeKeyUseCase: authorizeKeyUseCase,
	}
}
