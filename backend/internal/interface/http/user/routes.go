package user

import (
	"example.com/m/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.GET("/login", h.Login)
	e.GET("/callback", h.Callback)
	e.GET("/me", h.GetMe, middleware.AuthMiddleware(h.ensureUserUseCase))
	e.GET("/me/instances", h.ListMine, middleware.AuthMiddleware(h.ensureUserUseCase))
}
