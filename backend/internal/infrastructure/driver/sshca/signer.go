package sshca

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"example.com/m/internal/interface/http/access"
	"golang.org/x/crypto/ssh"
)

// https://pkg.go.dev/golang.org/x/crypto/ssh#Signer

type signer struct {
	caSigner ssh.Signer
	jumpUser string
}

func NewSigner(cfg *Config) access.Signer {
	keyBytes, err := os.ReadFile(cfg.privateKeyPath)
	if err != nil {
		panic(err)
	}

	caSigner, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		panic(err)
	}
	return &signer{
		caSigner: caSigner,
		jumpUser: cfg.jumpUser,
	}
}

func (s *signer) SignCertificate(pubKeyBytes []byte) (string, error) {
	// https://pkg.go.dev/golang.org/x/crypto/ssh#ParseAuthorizedKey
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	// ログを追える用にIDを付与する
	now := time.Now()
	keyID := fmt.Sprintf("%s-%d", s.jumpUser, now.Unix())
	// https://pkg.go.dev/golang.org/x/crypto/ssh#Certificate
	cert := &ssh.Certificate{
		Key:             pubKey,
		CertType:        ssh.UserCert,
		KeyId:           keyID,
		ValidPrincipals: []string{s.jumpUser},
		// 2時間前から24時間後まで有効な証明書を発行する
		// 期間設定については、のちに検討したいが、どうせ権限のないユーザーが踏み台なので、
		// 利便性を損なわない長めの期間にするのでよさそう
		ValidAfter:      uint64(now.Add(-2 * time.Hour).Unix()),
		ValidBefore:     uint64(now.Add(24 * time.Hour).Unix()),
	}

	if err := cert.SignCert(rand.Reader, s.caSigner); err != nil {
		return "", err
	}
	
	return string(ssh.MarshalAuthorizedKey(cert)), nil
}
