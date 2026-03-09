package network

import "context"

type NetworkService interface {
	CalculateNextAvailableIP(ctx context.Context, cidr string, usedIPs []string) (string, error)
}
