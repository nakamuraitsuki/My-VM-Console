package network

import (
	"context"
	"fmt"

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
	// VPCはIncusのProjectに対応させる
	put := api.ProjectPut{
		Description: fmt.Sprintf("VCP Project for %s", vpc.ID()),
		Config: map[string]string{
			"features.images":   "true",
			"features.networks": "true",
			"features.profiles": "true",
		},
	}

	req := api.ProjectsPost{
		ProjectPut: put,
		Name:       string(vpc.ID()),
	}

	err := d.client.CreateProject(req)
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) CreateSubnet(ctx context.Context, vpcID network.VPCID, subnet *network.Subnet) error {
	bridgeName := network.IDToResourceName(string(subnet.ID()))

	// project でVPCを表現
	pc := d.client.UseProject(string(vpcID))

	req := api.NetworksPost{
		Name: bridgeName,
		Type: "bridge",
		NetworkPut: api.NetworkPut{
			Config: map[string]string{
				"ipv4.address": subnet.CIDR(),
				"ipv4.nat":     "true",
			},
		},
	}
	err := pc.CreateNetwork(req)
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
