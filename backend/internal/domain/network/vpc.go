package network

import "github.com/google/uuid"

type VPCID string

type VPC struct {
	id      VPCID
	ownerID string // user.UserIDなど
	name    string
	cidr    string // 例: "10.0.0.0/16"
}

func NewVPC(id VPCID, ownerID, name, cidr string) *VPC {
	return &VPC{
		id:      id,
		ownerID: ownerID,
		name:    name,
		cidr:    cidr,
	}
}

func NewVPCID() VPCID {
	return VPCID("vpc-" + uuid.New().String())
}

// --- Getter ---
func (v *VPC) ID() VPCID       { return v.id }
func (v *VPC) OwnerID() string { return v.ownerID }
func (v *VPC) Name() string    { return v.name }
func (v *VPC) CIDR() string    { return v.cidr }
