package oidc

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

type IDTokenVerifier struct {
	verifier *oidc.IDTokenVerifier
}

// issuerURL は IDPのURL
// clientID は　IDP側に登録されている自サービスのID
func NewIDTokenVerifier(cfg *OIDCConfig) *IDTokenVerifier {
	ctx := context.Background()
	// 公開鍵取得と管理
	keySet := oidc.NewRemoteKeySet(
		ctx,
		fmt.Sprintf("%s/.well-known/jwks.json", cfg.IssuerURL),
	)

	// 検証設定
	config := &oidc.Config{
		ClientID: cfg.ClientID,
	}

	verifier := oidc.NewVerifier(cfg.IssuerURL, keySet, config)

	return &IDTokenVerifier{
		verifier: verifier,
	}
}

type VerifiedClaims struct {
	Subject string // OIDCのsubクレーム
	Name    string // OIDCのnameクレーム
	Nonce  string // OIDCのnonceクレーム
}

func (v *IDTokenVerifier) Verify(ctx context.Context, rawIDToken string) (*VerifiedClaims, error) {
	token, err := v.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	var claims struct {
		Subject string `json:"sub"`
		Name    string `json:"name"`
		Nonce   string `json:"nonce"`
	}
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &VerifiedClaims{
		Subject: claims.Subject,
		Name:    claims.Name,
		Nonce:   claims.Nonce,
	}, nil
}
