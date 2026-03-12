package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	// RFC6749 4.1.1
	// state はCSRF対策として送ることが推奨される
	state, _ := h.generateRandomString(32)
	// openID Connect core 3.1.2.1
	// nonce はリプレイ攻撃対策として送ることが推奨される
	nonce, _ := h.generateRandomString(32)

	// 覚えておかないとつかえないので、Cookieに保存しておく
	h.setAuthCookie(c, "state", state)
	h.setAuthCookie(c, "nonce", nonce)

	// openid と read:basic_info を IDPに要求する
	authURL := fmt.Sprintf(
		"%s/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=openid read:basic_info&state=%s&nonce=%s",
		h.oidcConfig.IssuerURL,
		h.oidcConfig.ClientID,
		h.oidcConfig.RedirectURL,
		state,
		nonce,
	)

	return c.Redirect(http.StatusFound, authURL)
}

// Cookie セット用のヘルパー
func (h *Handler) setAuthCookie(c echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   300,                  // 5分程度で十分
		HttpOnly: true,                 // JSから読み取らせない
		Secure:   false,                // 開発環境なら false、本番なら true (Configから取れると理想的)
		SameSite: http.SameSiteLaxMode, // リダイレクト時にCookieを維持するため
	}
	c.SetCookie(cookie)
}

// state や nonce 等を作るためのヘルパ
func (h *Handler) generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
