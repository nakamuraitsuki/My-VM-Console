package user

import (
	"github.com/labstack/echo/v4"
)

type ListMineResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	SubnetID  string `json:"subnet_id"`
	PrivateIP string `json:"private_ip"`
}

// AuthMiddleware でユーザー情報が Context にセットされている前提で、自分のインスタンスの一覧を返すエンドポイント
func (h *Handler) ListMine(c echo.Context) error {
	instances, err := h.listMyInstanceUseCase.Execute(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, "インスタンスの取得に失敗しました")
	}

	res := make([]ListMineResponse, len(instances))
	for i, instance := range instances {
		res[i] = ListMineResponse{
			ID:        string(instance.ID()),
			Name:      instance.Name(),
			Status:    string(instance.Status()),
			SubnetID:  string(instance.SubnetID()),
			PrivateIP: string(instance.PrivateIP()),
		}
	}
	return c.JSON(200, res)
}
