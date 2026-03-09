package network

type Lease struct {
	SubnetID SubnetID
	// インスタンス自体や、インスタンスのNICなど、
	// IPアドレスを割り当てる対象のID
	TargetID  string
	IPAddress string
}
