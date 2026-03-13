package etcd

import (
	"time"

	"example.com/m/internal/infrastructure/env"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewConfig() *clientv3.Config {
	return &clientv3.Config{
		Endpoints: []string{
			env.GetString("ETCD_ENDPOINT", "localhost:2379"),
		},
		DialTimeout: 5*time.Second, // とりあえずハードコーディング
	}
}
