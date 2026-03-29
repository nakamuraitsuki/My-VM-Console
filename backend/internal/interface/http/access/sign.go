package access

import (
	"github.com/labstack/echo/v4"
)

// NOTE: 中間層の認証スキップのために導入しているので、Principalsなどは固定する。
//	サーバーセットアップ時に、ジャンプ用に権限を絞ったユーザーが作成されていることを期待する
type SignRequest struct {
	PublicKey string `json:"public_key"`
}

func (h *Handler) Sign(c echo.Context) error {
	var req SignRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "invalid request body")
	}

	cert, err := h.signer.SignCertificate([]byte(req.PublicKey))
	if err != nil {
		return echo.NewHTTPError(500, "failed to sign certificate")
	}

	var res struct {
		Certificate string `json:"certificate"`
	}

	res.Certificate = cert
	return c.JSON(200, res)
}
