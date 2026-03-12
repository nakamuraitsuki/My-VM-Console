package compute

import (
	"example.com/m/internal/domain/image"
	"example.com/m/internal/domain/network"
	"example.com/m/internal/usecase/compute"
	"github.com/labstack/echo/v4"
)

type CreateReq struct {
	Name     string  `json:"name"`
	ImageID  string  `json:"image_id"`
	VPCID    string  `json:"vpc_id"`
	SubnetID *string `json:"subnet_id,omitempty"`
	Cpu      int     `json:"cpu"`
	Memory   int     `json:"memory"`
}

type CreateResp struct {
	InstanceID string `json:"instance_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	SubnetID   string `json:"subnet_id"`
	PrivateIP  string `json:"private_ip"`
}

func (h *Handler) CreateInstance(c echo.Context) error {
	var req CreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "リクエストのパースに失敗しました")
	}

	var subnetID *network.SubnetID
	if req.SubnetID != nil {
		sid := network.SubnetID(*req.SubnetID)
		subnetID = &sid
	}

	result, err := h.reqCreateUseCase.Execute(c.Request().Context(), compute.RequestCreateInstanceInput{
		Name:     req.Name,
		ImageID:  image.ImageID(req.ImageID),
		VPCID:    network.VPCID(req.VPCID),
		SubnetID: subnetID,
		CPU:      req.Cpu,
		Memory:   req.Memory,
	})
	if err != nil {
		return echo.NewHTTPError(500, "インスタンスの作成に失敗しました")
	}

	return c.JSON(200, CreateResp{
		InstanceID: string(result.InstanceID),
		Name:       result.Name,
		Status:     string(result.Status),
		SubnetID:   string(result.SubnetID),
		PrivateIP:  result.PrivateIP,
	})
}
