package network

import "context"

type NetworkService interface {
	CalculateNextAvailableIP(ctx context.Context, subnetID SubnetID) (string, error)
}
