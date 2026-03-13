package compute

import (
	"context"
	"fmt"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/image"
	"example.com/m/internal/domain/network"
	incus "github.com/lxc/incus/v6/client"
	"github.com/lxc/incus/v6/shared/api"
)

type driver struct {
	client incus.InstanceServer
}

func NewDriver(c incus.InstanceServer) compute.ComputeDriver {
	return &driver{
		client: c,
	}
}

func (d *driver) Create(ctx context.Context, inst *compute.Instance, img *image.Image) error {
	// 共通ルール
	bridgeName := network.IDToResourceName(string(inst.SubnetID()))

	// instance put
	put := api.InstancePut{
		Profiles: []string{"default"},
		Config: map[string]string{
			"limits.cpu":    fmt.Sprintf("%d", inst.CPU()),
			"limits.memory": fmt.Sprintf("%dMiB", inst.MemoryMB()),
			"user.owner":    string(inst.OwnerID()),
		},
		Devices: map[string]map[string]string{
			"eth0": {
				"type":         "nic",
				"nictype":      "bridged",
				"parent":       bridgeName,
				"name":         "eth0",
				"ipv4.address": inst.PrivateIP(),
			},
			"root": {
				"type": "disk",
				"path": "/",
				"pool": "default",
			},
		},
	}

	// instance post
	post := api.InstancesPost{
		Name:        string(inst.ID()),
		InstancePut: put,
		Source: api.InstanceSource{
			Type:        "image",
			Fingerprint: img.Fingerprint(),
			Server:      img.ServerURL(), // 例: "https://images.linuxcontainers.org"
			Protocol:    img.Protocol(),  // 例: "simplestreams"
			Mode:        "pull",
		},
	}

	op, err := d.client.CreateInstance(post)
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	err = op.WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for instance creation: %w", err)
	}
	return nil
}

func (d *driver) Start(ctx context.Context, id compute.InstanceID) error {
	req := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1, // 起動完了まで待機
	}

	op, err := d.client.UpdateInstanceState(string(id), req, "")
	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	if err := op.WaitContext(ctx); err != nil {
		return fmt.Errorf("error while waiting for instance start: %w", err)
	}
	return nil
}

func (d *driver) Stop(ctx context.Context, id compute.InstanceID) error {
	req := api.InstanceStatePut{
		Action:  "stop",
		Timeout: 30, // 30秒の猶予を持ってクリーンシャットダウンを試みる
		Force:   false,
	}

	op, err := d.client.UpdateInstanceState(string(id), req, "")
	if err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	if err := op.WaitContext(ctx); err != nil {
		return fmt.Errorf("error while waiting for instance stop: %w", err)
	}
	return nil
}

// NOTE: usecaseなどの層で、停止状態を保証する
func (d *driver) Terminate(ctx context.Context, id compute.InstanceID) error {
	op, err := d.client.DeleteInstance(string(id))
	if err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	if err := op.WaitContext(ctx); err != nil {
		return fmt.Errorf("error while waiting for instance deletion: %w", err)
	}
	return nil
}

// 本当にわけがわからなくなったときのリカバリ用。物理とDBの整合性を取りに行く
func (d *driver) GetRealStatus(ctx context.Context, id compute.InstanceID) (compute.InstanceStatus, error) {
	state, _, err := d.client.GetInstanceState(string(id))
	if err != nil {
		return compute.StatusError, fmt.Errorf("failed to get instance state: %w", err)
	}

	// IncusのStatusCodeをドメインのステータスへマッピング
	switch state.StatusCode {
	case api.Running:
		return compute.StatusRunning, nil
	case api.Stopped:
		return compute.StatusStopped, nil
	case api.Starting:
		return compute.StatusStarting, nil
	case api.Stopping:
		return compute.StatusStopping, nil
	case api.Aborting, api.Error:
		return compute.StatusError, nil
	default:
		return compute.StatusPending, nil
	}
}
