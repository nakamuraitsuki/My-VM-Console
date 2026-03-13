package network

import "context"

type NetworkService interface {
	CalculateNextAvailableVPCCIDR(ctx context.Context, usedCidrs []string) (string, error)
	CalculateNextAvailableSubnet(ctx context.Context, vpcCidr string, usedCidrs []string) (string, error)
	CalculateNextAvailableIP(ctx context.Context, cidr string, usedIPs []string) (string, error)
}
