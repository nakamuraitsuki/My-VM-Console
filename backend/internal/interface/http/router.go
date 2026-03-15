package http

import (
	"example.com/m/internal/infrastructure/env"
	"example.com/m/internal/interface/http/compute"
	"example.com/m/internal/interface/http/user"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func InitRoutes(
	userHandler *user.Handler,
	computeHandler *compute.Handler,
) *echo.Echo {
	e := echo.New()

	// セッションミドルウェアのDefault設定
	secret := env.GetString("SESSION_SECRET", "secret")
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
	}
	e.Use(session.Middleware(store))


	// ユーザ関連のルートを登録
	userGroup := e.Group("api/users")
	userHandler.RegisterRoutes(userGroup)

	// コンピュート関連のルートを登録
	computeGroup := e.Group("api/computes")
	computeHandler.RegisterRoutes(computeGroup)

	return e
}
