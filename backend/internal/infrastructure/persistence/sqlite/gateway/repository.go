package gateway

import (
	"context"
	"database/sql"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/gateway"
	"github.com/jmoiron/sqlx"
)

type ingressRouteModel struct {
	ID         string `db:"id"`
	Subdomain  string `db:"subdomain"`
	PortName   string `db:"port_name"`
	TargetIP   string `db:"target_ip"`
	TargetPort int    `db:"target_port"`
	InstanceID string `db:"instance_id"`
	OwnerID    string `db:"owner_id"`
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) gateway.Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(ctx context.Context, id gateway.IngressID) (*gateway.IngressRoute, error) {
	const query = `SELECT id, subdomain, port_name, target_ip, target_port, instance_id, owner_id FROM ingress_routes WHERE id = ?`
	var m ingressRouteModel
	if err := r.db.GetContext(ctx, &m, query, string(id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *repository) FindByInstanceID(ctx context.Context, instanceID compute.InstanceID) ([]*gateway.IngressRoute, error) {
	const query = `SELECT id, subdomain, port_name, target_ip, target_port, instance_id, owner_id FROM ingress_routes WHERE instance_id = ?`
	var models []ingressRouteModel
	if err := r.db.SelectContext(ctx, &models, query, string(instanceID)); err != nil {
		return nil, err
	}
	return toEntities(models), nil
}

func (r *repository) FindByOwnerID(ctx context.Context, ownerID string) ([]*gateway.IngressRoute, error) {
	const query = `SELECT id, subdomain, port_name, target_ip, target_port, instance_id, owner_id FROM ingress_routes WHERE owner_id = ?`
	var models []ingressRouteModel
	if err := r.db.SelectContext(ctx, &models, query, ownerID); err != nil {
		return nil, err
	}
	return toEntities(models), nil
}

func (r *repository) Save(ctx context.Context, route *gateway.IngressRoute) error {
	const query = `
INSERT INTO ingress_routes (id, subdomain, port_name, target_ip, target_port, instance_id, owner_id)
VALUES (:id, :subdomain, :port_name, :target_ip, :target_port, :instance_id, :owner_id)
ON CONFLICT(id) DO UPDATE SET
subdomain   = :subdomain,
port_name   = :port_name,
target_ip   = :target_ip,
target_port = :target_port,
instance_id = :instance_id,
owner_id    = :owner_id
`
	m := ingressRouteModel{
		ID:         string(route.ID()),
		Subdomain:  route.Subdomain(),
		PortName:   route.PortName(),
		TargetIP:   route.TargetIP(),
		TargetPort: route.TargetPort(),
		InstanceID: string(route.InstanceID()),
		OwnerID:    route.OwnerID(),
	}
	_, err := r.db.NamedExecContext(ctx, query, m)
	return err
}

func (r *repository) DeleteBulk(ctx context.Context, ids []gateway.IngressID) error {
	if len(ids) == 0 {
		return nil
	}
	raw := make([]interface{}, len(ids))
	for i, id := range ids {
		raw[i] = string(id)
	}
	query, args, err := sqlx.In(`DELETE FROM ingress_routes WHERE id IN (?)`, raw)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func toEntity(m *ingressRouteModel) *gateway.IngressRoute {
	return gateway.NewIngressRoute(
		gateway.IngressID(m.ID),
		m.Subdomain,
		m.PortName,
		m.TargetIP,
		m.TargetPort,
		m.OwnerID,
		compute.InstanceID(m.InstanceID),
	)
}

func toEntities(models []ingressRouteModel) []*gateway.IngressRoute {
	result := make([]*gateway.IngressRoute, len(models))
	for i := range models {
		result[i] = toEntity(&models[i])
	}
	return result
}
