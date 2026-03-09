package network

import "context"

type NetworkDriver interface {
	CreatePhysicalNetwork(ctx context.Context, vpc *VPC, subnet *Subnet) error
	GetNetworkStatus(ctx context.Context, vpcID VPCID) (string, error)
	DeletePhysicalNetwork(ctx context.Context, vpcID VPCID) error
}
