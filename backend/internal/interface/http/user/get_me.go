package user

import (
	"example.com/m/internal/domain/user"
	"github.com/labstack/echo/v4"
)

type GetMeResponse struct {
	ID              string            `json:"id"`
	DisplayName     string            `json:"display_name"`
	ProfileImageURL string            `json:"profile_image_url"`
	Permissions     []user.Permission `json:"permissions"`
	Quota           QuotaModel        `json:"quota"`
	Status          user.UserStatus   `json:"status"`
	ErrorPhase      *user.FailedPhase `json:"error_phase,omitempty"`
}

type QuotaModel struct {
	MaxInstance int `json:"max_instance"`
	MaxCPU      int `json:"max_cpu"`
	MaxMemory   int `json:"max_memory"`
}

func (h *Handler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()

	usr, ok := user.FromContext(ctx)
	if !ok {
		return echo.NewHTTPError(401, "ユーザー情報が見つかりません。")
	}

	response := GetMeResponse{
		ID:              string(usr.ID()),
		DisplayName:     usr.DisplayName(),
		ProfileImageURL: usr.ProfileImageURL(),
		Permissions:     usr.Permissions(),
		Quota: QuotaModel{
			MaxInstance: usr.Quota().MaxInstance,
			MaxCPU:      usr.Quota().MaxCPU,
			MaxMemory:   usr.Quota().MaxMemory,
		},
		Status:     usr.Status(),
		ErrorPhase: usr.ErrPhase(),
	}

	return c.JSON(200, response)
}
