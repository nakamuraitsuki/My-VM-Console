package network

import "github.com/google/uuid"

type SubnetID string

type Subnet struct {
	id    SubnetID
	vpcID VPCID
	name  string
	cidr  string // 例: "10.0.1.0/24"
}

// --- Constructor ---

func NewSubnet(
	id SubnetID,
	vpcID VPCID,
	name string,
	cidr string,
) *Subnet {
	return &Subnet{
		id:    id,
		vpcID: vpcID,
		name:  name,
		cidr:  cidr,
	}
}

func NewSubnetID() SubnetID {
	return SubnetID("sn-" + uuid.New().String())
}

// --- Getters ---

func (s *Subnet) ID() SubnetID { return s.id }
func (s *Subnet) VPCID() VPCID { return s.vpcID }
func (s *Subnet) Name() string { return s.name }
func (s *Subnet) CIDR() string { return s.cidr }
