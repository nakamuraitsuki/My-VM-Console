package middleware

import (
	"net/http"

	domainUser "example.com/m/internal/domain/user"
	"example.com/m/internal/usecase/user"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware はセッションを確認し、Domain Serviceを通じてユーザー情報を同期します
func AuthMiddleware(
	ensureUser user.EnsureUserUseCase,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)

			// セッションから ID と AccessToken を取り出す
			userID, okID := sess.Values["user_id"].(string)
			_, okTok := sess.Values["access_token"].(string) // 使わないが、あることがわかればいい

			// どちらか欠けていれば未認証として扱う
			if !okID || !okTok || userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "認証セッションが見つかりません。再ログインしてください。")
			}

			ctx := c.Request().Context()
			// UseCaseを呼び出してユーザー情報を同期
			resultUser, err := ensureUser.Execute(ctx, user.EnsureUserInput{
				Sub: userID,
			})
			if err != nil {
				// ここでエラーが出る場合は、トークン失効やDBエラーなどが考えられる
				return echo.NewHTTPError(http.StatusUnauthorized, "ユーザー情報の同期に失敗しました。")
			}

			// Echoのコンテキストにドメインモデルをセット
			domainUser.WithContext(ctx, resultUser)

			return next(c)
		}
	}
}