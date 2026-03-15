package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"example.com/m/internal/usecase/user"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// GET /api/callback
func (h *Handler) Callback(c echo.Context) error {
	ctx := c.Request().Context()
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	code := c.QueryParam("code")
	stateFromURL := c.QueryParam("state")

	stateCookie, err := c.Cookie("state")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "State cookie not found")
	}
	nonceCookie, err := c.Cookie("nonce")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Nonce cookie not found")
	}

	// CSRF対策: クエリパラメータのstateとクッキーのstateが一致するか確認
	if stateFromURL != stateCookie.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid state parameter")
	}

	tokenResponse, err := h.exchangeCodeForToken(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to exchange code for token: %v", err))
	}

	// AccessToken を焼く
	// NOTE: EnsureUserUseCaseの中で、APIによる権限補完が入るため、先にAccessTokenをセッションに保存している
	sess.Values["access_token"] = tokenResponse.AccessToken

	// id_token を検証してクレームを取得
	claims, err := h.idTokenVerifier.Verify(ctx, tokenResponse.IDToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Failed to verify ID token: %v", err))
	}

	if claims.Nonce != nonceCookie.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid nonce")
	}

	// JIT Provisioning
	ensureUserInput := user.EnsureUserInput{
		Sub:   claims.Subject,
		Token: tokenResponse.AccessToken,
	}
	resultUser, err := h.ensureUserUseCase.Execute(ctx, ensureUserInput)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to register user: %v", err))
	}

	// ユーザーIDを焼く
	sess.Values["user_id"] = string(resultUser.ID())
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	// クッキーをクリア
	h.clearTemporaryCookie(c, "state")
	h.clearTemporaryCookie(c, "nonce")

	return c.Redirect(http.StatusFound, "/dashboard")
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
}

func (h *Handler) exchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error) {
	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("redirect_uri", h.oidcConfig.RedirectURL)

	tokenEndpoint := fmt.Sprintf("%s/oauth/access-token", h.oidcConfig.IssuerURL)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenEndpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(h.oidcConfig.ClientID, h.oidcConfig.ClientSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	// トークンを解析して返す
	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &TokenResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		IDToken:      tokenResponse.IDToken,
	}, nil
}

// 一時的に使用したクッキーを削除するためのヘルパー関数
func (h *Handler) clearTemporaryCookie(c echo.Context, name string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // これでブラウザは即座にCookieを破棄する
		HttpOnly: true,
	})
}
