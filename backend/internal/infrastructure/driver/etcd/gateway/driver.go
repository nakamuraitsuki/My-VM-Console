package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"example.com/m/internal/domain/gateway"
	"example.com/m/internal/infrastructure/env"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type driver struct {
	client     *clientv3.Client
	baseDomain string
	prefix     string
}

func NewDriver(c *clientv3.Client) gateway.IngressDriver {
	return &driver{
		client:     c,
		baseDomain: env.GetString("BASE_DOMAIN", "localhost"),
		prefix:     env.GetString("ETCD_INGRESS_PREFIX", "/skydns"),
	}
}

func (d *driver) ApplyRoutes(ctx context.Context, routes []*gateway.IngressRoute) error {
	for _, route := range routes {
		fullDomain := d.buildFullDomain(route.Subdomain())
		path := d.domainToEtcdPath(fullDomain)

		// record を etcd に保存
		record := map[string]any{
			"host": route.TargetIP(),
			"port": route.TargetPort(),
			"ttl":  10, // TTLはとりあえず固定値
		}

		val, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %w", err)
		}

		_, err = d.client.Put(ctx, path, string(val))
		if err != nil {
			return fmt.Errorf("failed to put record to etcd: %w", err)
		}
	}
	return nil
}

func (d *driver) RemoveRoutes(ctx context.Context, routes []*gateway.IngressRoute) error {
	for _, route := range routes {
		fullDomain := d.buildFullDomain(route.Subdomain())
		path := d.domainToEtcdPath(fullDomain)

		_, err := d.client.Delete(ctx, path)
		if err != nil {
			return fmt.Errorf("failed to delete record from etcd: %w", err)
		}
	}
	return nil
}

// --- helper ---
func (d *driver) buildFullDomain(subdomain string) string {
	return fmt.Sprintf("%s.%s", subdomain, d.baseDomain)
}

func (d *driver) domainToEtcdPath(fullDomain string) string {
	parts := strings.Split(fullDomain, ".")

	// 親方向から解決したいので、partsを逆順にする
	slices.Reverse(parts)
	return fmt.Sprintf("%s/%s", d.prefix, strings.Join(parts, "/"))
}
