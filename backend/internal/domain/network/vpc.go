package network

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
