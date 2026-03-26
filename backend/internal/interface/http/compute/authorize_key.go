package compute

import (
	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/user"
	computeUC "example.com/m/internal/usecase/compute"
	"github.com/labstack/echo/v4"
)

type AuthorizePubKeyRequest struct {
	InstanceID string `json:"instance_id"`
	PubKey     string `json:"pub_key"`
}

func (h *Handler) AuthorizePubKey(c echo.Context) error {
	var req AuthorizePubKeyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "invalid request body")
	}

	input := computeUC.AuthorizePubKeyInput{
		InstanceID: compute.InstanceID(req.InstanceID),
		PubKey:     req.PubKey,
	}

	if err := h.authorizeKeyUseCase.Execute(c.Request().Context(), input); err != nil {
		if err == compute.ErrInstanceNotFound {
			return echo.NewHTTPError(404, "instance not found")
		}
		if err == user.ErrUserNotInContext {
			return echo.NewHTTPError(401, "user not authenticated")
		}
		if err == user.ErrNoPermission {
			return echo.NewHTTPError(403, "forbidden")
		}
		if err == compute.ErrInvalidInstanceStatus {
			return echo.NewHTTPError(400, "invalid instance status")
		}
		return echo.NewHTTPError(500, "internal server error")
	}

	return c.NoContent(204)
}
