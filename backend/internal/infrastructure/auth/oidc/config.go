package oidc

import "example.com/m/internal/infrastructure/env"

type OIDCConfig struct {
	IssuerURL string
	ClientID string
	ClientSecret string
	RedirectURL string
}

func NewOIDCConfig() *OIDCConfig {
	return &OIDCConfig{
		IssuerURL: env.GetString("OIDC_ISSUER_URL", "http://oauth.local:8080"),
		ClientID: env.GetString("OIDC_CLIENT_ID", "a394037d-727b-499e-b9b4-3a78c7615fef"),
		ClientSecret: env.GetString("OIDC_CLIENT_SECRET", "secret"),
		RedirectURL: env.GetString("OIDC_REDIRECT_URL", "http://localhost/api/callback"),
	}
}