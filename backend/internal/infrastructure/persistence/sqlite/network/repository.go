package network

import (
	"context"
	"database/sql"

	"example.com/m/internal/domain/network"
	"github.com/jmoiron/sqlx"
)

type vpcModel struct {
	ID      string `db:"id"`
	OwnerID string `db:"owner_id"`
	Name    string `db:"name"`
	CIDR    string `db:"cidr"`
}

type subnetModel struct {
	ID    string `db:"id"`
	VPCID string `db:"vpc_id"`
	Name  string `db:"name"`
	CIDR  string `db:"cidr"`
}

type leaseModel struct {
	SubnetID  string `db:"subnet_id"`
	TargetID  string `db:"target_id"`
	IPAddress string `db:"ip_address"`
}

type NetworkRepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) network.Repository {
	return &NetworkRepositoryImpl{db: db}
}

func (r *NetworkRepositoryImpl) FindVPCByID(ctx context.Context, id network.VPCID) (*network.VPC, error) {
	const query = `SELECT id, owner_id, name, cidr FROM vpcs WHERE id = ?`
	var m vpcModel
	if err := r.db.GetContext(ctx, &m, query, string(id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return network.NewVPC(network.VPCID(m.ID), m.OwnerID, m.Name, m.CIDR), nil
}

func (r *NetworkRepositoryImpl) FindSubnetsByVPCID(ctx context.Context, vpcID network.VPCID) ([]*network.Subnet, error) {
	const query = `SELECT id, vpc_id, name, cidr FROM subnets WHERE vpc_id = ?`
	var models []subnetModel
	if err := r.db.SelectContext(ctx, &models, query, string(vpcID)); err != nil {
		return nil, err
	}

	result := make([]*network.Subnet, len(models))
	for i, m := range models {
		result[i] = network.NewSubnet(
			network.SubnetID(m.ID),
			network.VPCID(m.VPCID),
			m.Name,
			m.CIDR,
		)
	}
	return result, nil
}

func (r *NetworkRepositoryImpl) FindSubnetByID(ctx context.Context, id network.SubnetID) (*network.Subnet, error) {
	const query = `SELECT id, vpc_id, name, cidr FROM subnets WHERE id = ?`
	var m subnetModel
	if err := r.db.GetContext(ctx, &m, query, string(id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return network.NewSubnet(network.SubnetID(m.ID), network.VPCID(m.VPCID), m.Name, m.CIDR), nil
}

func (r *NetworkRepositoryImpl) SaveVPC(ctx context.Context, vpc *network.VPC) error {
	model := vpcModel{
		ID:      string(vpc.ID()),
		OwnerID: vpc.OwnerID(),
		Name:    vpc.Name(),
		CIDR:    vpc.CIDR(),
	}

	const query = `
INSERT INTO vpcs (id, owner_id, name, cidr)
VALUES (:id, :owner_id, :name, :cidr)
ON CONFLICT(id) DO UPDATE SET
	owner_id = :owner_id,
	name = :name,
	cidr = :cidr
`

	_, err := r.db.NamedExecContext(ctx, query, model)
	return err
}

func (r *NetworkRepositoryImpl) SaveSubnet(ctx context.Context, subnet *network.Subnet) error {
	model := subnetModel{
		ID:    string(subnet.ID()),
		VPCID: string(subnet.VPCID()),
		Name:  subnet.Name(),
		CIDR:  subnet.CIDR(),
	}

	const query = `
INSERT INTO subnets (id, vpc_id, name, cidr)
VALUES (:id, :vpc_id, :name, :cidr)
ON CONFLICT(id) DO UPDATE SET
	vpc_id = :vpc_id,
	name = :name,
	cidr = :cidr
`

	_, err := r.db.NamedExecContext(ctx, query, model)
	return err
}

func (r *NetworkRepositoryImpl) DeleteVPC(ctx context.Context, id network.VPCID) error {
	const query = `DELETE FROM vpcs WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, string(id))
	return err
}

func (r *NetworkRepositoryImpl) DeleteSubnet(ctx context.Context, id network.SubnetID) error {
	const query = `DELETE FROM subnets WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, string(id))
	return err
}

func (r *NetworkRepositoryImpl) ListAllUsedCIDRs(ctx context.Context) ([]string, error) {
	const query = `SELECT cidr FROM vpcs`
	var cidrs []string
	if err := r.db.SelectContext(ctx, &cidrs, query); err != nil {
		return nil, err
	}
	return cidrs, nil
}

func (r *NetworkRepositoryImpl) FindLeaseByTargetID(ctx context.Context, targetID string) (*network.Lease, error) {
	const query = `SELECT subnet_id, target_id, ip_address FROM leases WHERE target_id = ?`
	var m leaseModel
	if err := r.db.GetContext(ctx, &m, query, targetID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return network.NewLease(
		network.SubnetID(m.SubnetID),
		m.TargetID,
		m.IPAddress,
	), nil
}

func (r *NetworkRepositoryImpl) FindLeasesBySubnetID(ctx context.Context, subnetID network.SubnetID) ([]*network.Lease, error) {
	const query = `SELECT subnet_id, target_id, ip_address FROM leases WHERE subnet_id = ?`
	var models []leaseModel
	if err := r.db.SelectContext(ctx, &models, query, string(subnetID)); err != nil {
		return nil, err
	}

	result := make([]*network.Lease, len(models))
	for i, m := range models {
		result[i] = network.NewLease(
			network.SubnetID(m.SubnetID),
			m.TargetID,
			m.IPAddress,
		)
	}
	return result, nil
}

func (r *NetworkRepositoryImpl) CreateLease(ctx context.Context, lease *network.Lease) error {
	model := leaseModel{
		SubnetID:  string(lease.SubnetID),
		TargetID:  lease.TargetID,
		IPAddress: lease.IPAddress,
	}

	const query = `
INSERT INTO leases (subnet_id, target_id, ip_address)
VALUES (:subnet_id, :target_id, :ip_address)
ON CONFLICT(target_id) DO UPDATE SET
	subnet_id = :subnet_id,
	ip_address = :ip_address
`

	_, err := r.db.NamedExecContext(ctx, query, model)
	return err
}

func (r *NetworkRepositoryImpl) DeleteLease(ctx context.Context, targetID string) error {
	const query = `DELETE FROM leases WHERE target_id = ?`
	_, err := r.db.ExecContext(ctx, query, targetID)
	return err
}
