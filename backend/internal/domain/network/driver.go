package network

import "context"

type NetworkDriver interface {
	// VPC (Bridgeの作成)
	CreateVPC(ctx context.Context, vpc *VPC) error

	// Subnet (BridgeへのIP割り当てやDHCP設定)
	CreateSubnet(ctx context.Context, vpcID VPCID, subnet *Subnet) error

	// 削除も段階的に行えるようにする
	DeleteSubnet(ctx context.Context, subnetID SubnetID) error
	DeleteVPC(ctx context.Context, vpcID VPCID) error

	// 疎通確認用
	IsVPCReady(ctx context.Context, vpcID VPCID) (bool, error)
}
