package network

type Lease struct {
	SubnetID SubnetID
	// インスタンス自体や、インスタンスのNICなど、
	// IPアドレスを割り当てる対象のID
	TargetID  string
	IPAddress string
}

// --- Constructor ---

func NewLease(
	subnetID SubnetID,
	targetID string,
	ipAddress string,
) *Lease {
	return &Lease{
		SubnetID:  subnetID,
		TargetID:  targetID,
		IPAddress: ipAddress,
	}
}
