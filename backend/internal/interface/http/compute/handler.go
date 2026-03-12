package compute

import (
	"example.com/m/internal/usecase/compute"
	"example.com/m/internal/usecase/user"
)

type Handler struct {
	reqCreateUseCase  compute.RequestCreateInstanceUseCase
	ensureUserUseCase user.EnsureUserUseCase
}

func NewHandler(
	reqCreateUseCase compute.RequestCreateInstanceUseCase,
	ensureUserUseCase user.EnsureUserUseCase,
) *Handler {
	return &Handler{
		reqCreateUseCase:  reqCreateUseCase,
		ensureUserUseCase: ensureUserUseCase,
	}
}
