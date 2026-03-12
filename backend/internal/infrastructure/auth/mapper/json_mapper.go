package mapper

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"example.com/m/internal/domain/user"
)

type permissionMapper struct {
	mu       sync.RWMutex
	roleMap  map[string][]user.Permission
	filePath string
}

func NewPermissionMapper(filePath string, pollInterval time.Duration) user.PermissionMapper {
	m := &permissionMapper{
		filePath: filePath,
		roleMap:  make(map[string][]user.Permission),
	}

	// 起動時に一回読み込む
	if err := m.reload(); err != nil {
		log.Printf("[Mapper] Initial load failed: %v. Starting with empty map.", err)
	}

	// 定期ポーリング開始
	go m.startPolling(pollInterval)

	return m
}

func (m *permissionMapper) startPolling(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.reload(); err != nil {
			log.Printf("[Mapper] Periodic reload failed: %v. Keeping previous configuration.", err)
		}
	}
}

func (m *permissionMapper) reload() error {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	var config struct {
		Mappings map[string][]string `json:"mappings"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// 新しいマップを一時的に作成（文字列からPermission定数への変換）
	newMap := make(map[string][]user.Permission)
	for role, perms := range config.Mappings {
		for _, pStr := range perms {
			// pStr を Permission 型にキャスト、またはバリデーション
			newMap[role] = append(newMap[role], user.Permission(pStr))
		}
	}

	// ロックをかけて差し替え
	m.mu.Lock()
	m.roleMap = newMap
	m.mu.Unlock()

	return nil
}

func (m *permissionMapper) MapToPermissions(roleIDs []string) []user.Permission {
	m.mu.RLock()
	defer m.mu.RUnlock()

	uniquePerms := make(map[user.Permission]struct{})
	for _, roleID := range roleIDs {
		if perms, ok := m.roleMap[roleID]; ok {
			for _, p := range perms {
				uniquePerms[p] = struct{}{}
			}
		}
	}

	// スライス化
	res := make([]user.Permission, 0, len(uniquePerms))
	for p := range uniquePerms {
		res = append(res, p)
	}
	return res
}
