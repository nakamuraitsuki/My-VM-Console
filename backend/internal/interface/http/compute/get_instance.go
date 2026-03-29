package compute

import (
	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/user"
	computeUC "example.com/m/internal/usecase/compute"
	"github.com/labstack/echo/v4"
)

type GetInstanceResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	SubnetID   string   `json:"subnet_id"`
	PrivateIP  string   `json:"private_ip"`
	Subdomains []string `json:"subdomains"`
}

func (h *Handler) GetInstance(c echo.Context) error {
	instanceID := compute.InstanceID(c.Param("instance_id"))

	output, err := h.getInstanceUseCase.Execute(c.Request().Context(), computeUC.GetInstanceInput{
		InstanceID: instanceID,
	})
	if err != nil {
		if err == compute.ErrInstanceNotFound {
			return echo.NewHTTPError(404, "instance not found")
		}
		if err == user.ErrUserNotInContext {
			return echo.NewHTTPError(401, "user not authenticated")
		}
		if err == user.ErrNoPermission {
			return echo.NewHTTPError(403, "forbidden")
		}
		return echo.NewHTTPError(500, "internal server error")
	}

	subdomains := make([]string, 0, len(output.Ingresses))
	for _, ingress := range output.Ingresses {
		subdomains = append(subdomains, ingress.Subdomain())
	}

	inst := output.Instance
	return c.JSON(200, GetInstanceResponse{
		ID:         string(inst.ID()),
		Name:       inst.Name(),
		Status:     string(inst.Status()),
		SubnetID:   string(inst.SubnetID()),
		PrivateIP:  inst.PrivateIP(),
		Subdomains: subdomains,
	})
}
