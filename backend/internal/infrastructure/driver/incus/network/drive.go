package network

import (
	"context"
	"fmt"
	"net/netip"

	"example.com/m/internal/domain/network"
	incus "github.com/lxc/incus/v6/client"
	"github.com/lxc/incus/v6/shared/api"
)

type driver struct {
	client incus.InstanceServer
}

func NewDriver(c incus.InstanceServer) network.NetworkDriver {
	return &driver{
		client: c,
	}
}

func (d *driver) CreateVPC(ctx context.Context, vpc *network.VPC) error {
	vpcName := network.IDToResourceName(string(vpc.ID()))

	// VPCはIncusのProjectに対応させる
	put := api.ProjectPut{
		Description: fmt.Sprintf("VCP Project for %s", vpc.ID()),
		Config: map[string]string{
			"features.images":             "true",
			"features.networks":           "true",
			"features.profiles":           "true",
			"restricted":                  "true",
			"restricted.networks.uplinks": fmt.Sprintf("ovn-uplink, %s", vpcName),
		},
	}

	req := api.ProjectsPost{
		ProjectPut: put,
		Name:       string(vpc.ID()), // Project name
	}

	err := d.client.CreateProject(req)
	if err != nil {
		return err
	}

	// VPC専用のネットワーク生成
	pc := d.client.UseProject(string(vpc.ID()))

	prefix, err := netip.ParsePrefix(vpc.CIDR())
	if err != nil {
		return fmt.Errorf("invalid VPC CIDR: %w", err)
	}
	vpcListenAddr := fmt.Sprintf("%s/%d", prefix.Addr().Next().String(), prefix.Bits())

	mainNetwork := api.NetworksPost{
		Name: vpcName, // Network name
		Type: "ovn",
		NetworkPut: api.NetworkPut{
			Config: map[string]string{
				"network":      "ovn-uplink",
				"ipv4.address": vpcListenAddr,
				"ipv4.nat":     "false",
			},
		},
	}

	err = pc.CreateNetwork(mainNetwork)
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) CreateSubnet(ctx context.Context, vpcID network.VPCID, subnet *network.Subnet) error {
	bridgeName := network.IDToResourceName(string(subnet.ID()))
	vpcName := network.IDToResourceName(string(vpcID))

	// NOTE: subnet は ネットワークアドレスが渡されるので、Gateway設定は+1してあててあげる
	prefix, err := netip.ParsePrefix(subnet.CIDR())
	if err != nil {
		return fmt.Errorf("invalid subnet CIDR: %w", err)
	}
	// Incusの ipv4.address は「そのネットワーク内でIncus自身が持つIP」を期待するため
	gatewayAddr := prefix.Addr().Next()
	listenAddr := fmt.Sprintf("%s/%d", gatewayAddr.String(), prefix.Bits())

	// project でVPCを表現
	pc := d.client.UseProject(string(vpcID))

	req := api.NetworksPost{
		Name: bridgeName,
		Type: "ovn",
		NetworkPut: api.NetworkPut{
			Config: map[string]string{
				"ipv4.address":     listenAddr,
				"ipv4.nat":         "false",
				"ipv4.dhcp":        "true",
				"ipv4.dhcp.ranges": "", // IP割り振りはすべてGoアプリ側で行うが、その他の情報はほしい
			},
		},
	}
	err = pc.CreateNetwork(req)
	if err != nil {
		return err
	}

	// 双方向の Peering を設定
	err = pc.CreateNetworkPeer(bridgeName, api.NetworkPeersPost{
		Name:          "to-vpc-hub",
		TargetProject: string(vpcID),
		TargetNetwork: vpcName,
	})
	if err != nil {
		return err
	}

	err = pc.CreateNetworkPeer(vpcName, api.NetworkPeersPost{
		Name:          "to-" + bridgeName,
		TargetProject: string(vpcID),
		TargetNetwork: bridgeName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) DeleteSubnet(ctx context.Context, vpcID network.VPCID, subnetID network.SubnetID) error {
	bridgeName := network.IDToResourceName(string(subnetID))
	pc := d.client.UseProject(string(vpcID))
	return pc.DeleteNetwork(bridgeName)
}

func (d *driver) DeleteVPC(ctx context.Context, vpcID network.VPCID) error {
	return d.client.DeleteProject(string(vpcID))
}

func (d *driver) IsVPCReady(ctx context.Context, vpcID network.VPCID) (bool, error) {
	project, _, err := d.client.GetProject(string(vpcID))
	if err != nil {
		return false, err
	}

	// プロジェクトが存在すればいったんOK
	return project.Name == string(vpcID), nil
}
