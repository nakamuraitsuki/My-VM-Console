package network

import (
	"context"
	"database/sql"

	"example.com/m/internal/domain/network"
	"example.com/m/internal/infrastructure/persistence/sqlite"
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

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) network.Repository {
	return &repository{db: db}
}

func (r *repository) FindVPCByID(ctx context.Context, id network.VPCID) (*network.VPC, error) {
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

func (r *repository) FindSubnetsByVPCID(ctx context.Context, vpcID network.VPCID) ([]*network.Subnet, error) {
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

func (r *repository) FindSubnetByID(ctx context.Context, id network.SubnetID) (*network.Subnet, error) {
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

func (r *repository) SaveVPC(ctx context.Context, vpc *network.VPC) error {
	db := sqlite.GetExt(ctx, r.db)

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

	_, err := sqlx.NamedExecContext(ctx, db, query, model)
	return err
}

func (r *repository) SaveSubnet(ctx context.Context, subnet *network.Subnet) error {
	db := sqlite.GetExt(ctx, r.db)

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

	_, err := sqlx.NamedExecContext(ctx, db, query, model)
	return err
}

func (r *repository) DeleteVPC(ctx context.Context, id network.VPCID) error {
	db := sqlite.GetExt(ctx, r.db)
	const query = `DELETE FROM vpcs WHERE id = ?`
	_, err := db.ExecContext(ctx, query, string(id))
	return err
}

func (r *repository) DeleteSubnet(ctx context.Context, id network.SubnetID) error {
	db := sqlite.GetExt(ctx, r.db)
	const query = `DELETE FROM subnets WHERE id = ?`
	_, err := db.ExecContext(ctx, query, string(id))
	return err
}

func (r *repository) ListAllUsedCIDRs(ctx context.Context) ([]string, error) {
	const query = `SELECT cidr FROM vpcs`
	var cidrs []string
	if err := r.db.SelectContext(ctx, &cidrs, query); err != nil {
		return nil, err
	}
	return cidrs, nil
}

func (r *repository) FindLeaseByTargetID(ctx context.Context, targetID string) (*network.Lease, error) {
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

func (r *repository) FindLeasesBySubnetID(ctx context.Context, subnetID network.SubnetID) ([]*network.Lease, error) {
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

func (r *repository) CreateLease(ctx context.Context, lease *network.Lease) error {
	db := sqlite.GetExt(ctx, r.db)

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

	_, err := sqlx.NamedExecContext(ctx, db, query, model)
	return err
}

func (r *repository) DeleteLease(ctx context.Context, targetID string) error {
	db := sqlite.GetExt(ctx, r.db)
	const query = `DELETE FROM leases WHERE target_id = ?`
	_, err := db.ExecContext(ctx, query, targetID)
	return err
}
