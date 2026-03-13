package user

import (
	"github.com/labstack/echo/v4"
)

// AuthMiddleware でユーザー情報が Context にセットされている前提で、自分のインスタンスの一覧を返すエンドポイント
func (h *Handler) ListMine(c echo.Context) error {
	instances, err := h.listMyInstanceUseCase.Execute(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, "インスタンスの取得に失敗しました")
	}

	return c.JSON(200, instances)
}
