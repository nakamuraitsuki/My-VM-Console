package compute

import (
	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.POST("/instances", h.CreateInstance, middleware.AuthMiddleware(h.ensureUserUseCase))
}
