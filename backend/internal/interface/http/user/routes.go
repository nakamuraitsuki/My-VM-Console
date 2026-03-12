package user

import "github.com/labstack/echo/v4"

func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.GET("/login", h.Login)
	e.GET("/callback", h.Callback)
}
