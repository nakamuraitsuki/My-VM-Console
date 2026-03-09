package network

type SubnetID string

type Subnet struct {
	id    SubnetID
	vpcID VPCID
	name  string
	cidr  string // 例: "10.0.1.0/24"
}
