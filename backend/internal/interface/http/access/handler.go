package access

import "example.com/m/internal/usecase/user"

// NOTE: 認証スキップのために導入しているので、ユーザーIDなどはInfra層で固定することを想定している
type Signer interface {
	SignCertificate(pubKey []byte) (string, error)
}

type Handler struct {
	signer Signer
	ensureUserUseCase user.EnsureUserUseCase
}

func NewHandler(signer Signer, ensureUserUseCase user.EnsureUserUseCase) *Handler {
	return &Handler{signer: signer, ensureUserUseCase: ensureUserUseCase}
}