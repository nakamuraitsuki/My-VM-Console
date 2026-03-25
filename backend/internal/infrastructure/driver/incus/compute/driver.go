package compute

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

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
	vpcClient := d.client.UseProject(string(inst.VPCID()))

	// Defaultのユーザーを作成しておく。
	userData := `#cloud-config
users:
  - name: default
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
`

	// instance put
	put := api.InstancePut{
		Profiles: []string{"default"},
		Config: map[string]string{
			"limits.cpu":     fmt.Sprintf("%d", inst.CPU()),
			"limits.memory":  fmt.Sprintf("%dMiB", inst.MemoryMB()),
			"user.user-data": userData,
			"user.owner":     string(inst.OwnerID()),
		},
		Devices: map[string]map[string]string{
			"eth0": {
				"type":         "nic",
				"nictype":      "ovn",
				"network":      bridgeName,
				"name":         "eth0",
				"ipv4.address": inst.PrivateIP(),
			},
			"root": {
				"type": "disk",
				"path": "/",
				"pool": "default",
				// "source": string(inst.RootVolumeID()), // Volumeは自動解決されるので、競合を防ぎたい
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
			Alias:       img.Alias(),
			Server:      img.ServerURL(), // 例: "https://images.linuxcontainers.org"
			Protocol:    img.Protocol(),  // 例: "simplestreams"
			Mode:        "pull",
		},
	}

	op, err := vpcClient.CreateInstance(post)
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	err = op.WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for instance creation: %w", err)
	}
	return nil
}

func (d *driver) Start(ctx context.Context, inst *compute.Instance) error {
	vpcClient := d.client.UseProject(string(inst.VPCID()))
	req := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1, // 起動完了まで待機
	}

	op, err := vpcClient.UpdateInstanceState(string(inst.ID()), req, "")
	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	if err := op.WaitContext(ctx); err != nil {
		return fmt.Errorf("error while waiting for instance start: %w", err)
	}
	return nil
}

func (d *driver) Stop(ctx context.Context, inst *compute.Instance) error {
	vpcClient := d.client.UseProject(string(inst.VPCID()))
	req := api.InstanceStatePut{
		Action:  "stop",
		Timeout: 30, // 30秒の猶予を持ってクリーンシャットダウンを試みる
		Force:   false,
	}

	op, err := vpcClient.UpdateInstanceState(string(inst.ID()), req, "")
	if err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	if err := op.WaitContext(ctx); err != nil {
		return fmt.Errorf("error while waiting for instance stop: %w", err)
	}
	return nil
}

// NOTE: usecaseなどの層で、停止状態を保証する
func (d *driver) Terminate(ctx context.Context, inst *compute.Instance) error {
	vpcClient := d.client.UseProject(string(inst.VPCID()))
	op, err := vpcClient.DeleteInstance(string(inst.ID()))
	if err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	if err := op.WaitContext(ctx); err != nil {
		return fmt.Errorf("error while waiting for instance deletion: %w", err)
	}
	return nil
}

func (d *driver) AuthorizePublicKey(ctx context.Context, inst *compute.Instance, publicKey string) error {
	vpcClient := d.client.UseProject(string(inst.VPCID()))
	instanceName := string(inst.ID())
	path := "/home/default/.ssh/authorized_keys"

	var currentKeys string

	// 1. 現在の authorized_keys を取得
	content, _, err := vpcClient.GetInstanceFile(instanceName, path)
	if err != nil {
		return err
	}
	defer content.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, content)
	if err != nil {
		// ファイルが存在しない（ディレクトリがまだない等）場合は空として扱う
		currentKeys = ""
	}

	// 冪等性のチェック
	trimmedKey := strings.TrimSpace(publicKey)
	if strings.Contains(currentKeys, trimmedKey) {
		return nil // 既に存在すれば何もしない
	}

	// 追記
	var sb strings.Builder
	// 既存キー丸うつし
	sb.WriteString(strings.TrimRight(currentKeys, "\n"))
	if sb.Len() > 0 {
		sb.WriteString("\n")
	}
	// 新しいキーを追加
	sb.WriteString(trimmedKey)
	sb.WriteString("\n")

	// 書き戻し
	args := incus.InstanceFileArgs{
		Content:   strings.NewReader(sb.String()),
		Mode:      0600, // defaultユーザーが所有者であることを前提としたパーミッション
		Type:      "file",
		WriteMode: "overwrite",
	}
	// UID/GIDをdefaultユーザー(通常1000)に合わせる必要があるかも
	args.UID = 1000
	args.GID = 1000

	return vpcClient.CreateInstanceFile(instanceName, path, args)
}

// TODO: 実装
func (d *driver) RevokePublicKey(ctx context.Context, inst *compute.Instance, publicKey string) error {
	panic("not implemented")
}

// 本当にわけがわからなくなったときのリカバリ用。物理とDBの整合性を取りに行く
func (d *driver) GetRealStatus(ctx context.Context, inst *compute.Instance) (compute.InstanceStatus, error) {
	vpcClient := d.client.UseProject(string(inst.VPCID()))
	state, _, err := vpcClient.GetInstanceState(string(inst.ID()))
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
