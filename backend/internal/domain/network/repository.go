package network

import "context"

type Repository interface {
	FindVPCByID(ctx context.Context, id VPCID) (*VPC, error)
	FindSubnetsByVPCID(ctx context.Context, vpcID VPCID) ([]*Subnet, error)
	FindSubnetByID(ctx context.Context, id SubnetID) (*Subnet, error)
	SaveVPC(ctx context.Context, vpc *VPC) error
	SaveSubnet(ctx context.Context, subnet *Subnet) error
	DeleteVPC(ctx context.Context, id VPCID) error
	DeleteSubnet(ctx context.Context, id SubnetID) error

	FindLeaseByTargetID(ctx context.Context, targetID string) (*Lease, error)
	FindLeasesBySubnetID(ctx context.Context, subnetID SubnetID) ([]*Lease, error)
	CreateLease(ctx context.Context, lease *Lease) error
	DeleteLease(ctx context.Context, targetID string) error
}
