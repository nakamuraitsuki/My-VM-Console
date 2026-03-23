package access

import (
	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.POST("/sign", h.Sign, middleware.AuthMiddleware(h.ensureUserUseCase))
}
