package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"example.com/m/internal/domain/user"
	"github.com/patrickmn/go-cache"
)

type idpResponse struct {
	ID              string   `json:"id"`
	DisplayID       string   `json:"display_id"`
	DisplayName     string   `json:"display_name"`
	ProfileImageURL string   `json:"profile_image_url"`
	Roles           []string `json:"roles"`
}

type externalIdentityService struct {
	client      *http.Client
	apiEndpoint string
	mapper      user.PermissionMapper
	cache       *cache.Cache // キャッシュ用
}

func NewExternalIdentityService(client *http.Client, apiEndpoint string, mapper user.PermissionMapper) user.IdentityService {
	return &externalIdentityService{
		client:      client,
		apiEndpoint: apiEndpoint,
		mapper:      mapper,
		cache:       cache.New(30*time.Minute, 60*time.Minute), // キャッシュの有効期限とクリーンアップ間隔を設定
	}
}

func (s *externalIdentityService) GetIdentity(ctx context.Context, token string) (*user.UserIdentity, error) {
	// まずキャッシュを確認
	if cachedIdentity, found := s.cache.Get(token); found {
		if identity, ok := cachedIdentity.(*user.UserIdentity); ok {
			return identity, nil
		}
	}

	// キャッシュにない場合は外部APIを呼び出す
	req, err := http.NewRequestWithContext(ctx, "GET", s.apiEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user identity from external service")
	}

	var apiResponse idpResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	permissions := s.mapper.MapToPermissions(apiResponse.Roles)

	// パースしたデータをドメインモデルに変換
	identity := &user.UserIdentity{
		ID:              user.UserID(apiResponse.ID),
		DisplayName:     apiResponse.DisplayName,
		Permissions:     permissions,
	}

	// キャッシュに保存
	s.cache.Set(token, identity, cache.DefaultExpiration)

	return identity, nil
}
