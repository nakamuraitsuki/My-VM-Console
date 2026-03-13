package compute

import (
	"context"
	"database/sql"

	"example.com/m/internal/domain/compute"
	"example.com/m/internal/domain/image"
	"example.com/m/internal/domain/network"
	"example.com/m/internal/domain/storage"
	"example.com/m/internal/domain/user"
	"example.com/m/internal/infrastructure/persistence/sqlite"
	"github.com/jmoiron/sqlx"
)

type instanceModel struct {
	ID           string  `db:"id"`
	Name         string  `db:"name"`
	OwnerID      string  `db:"owner_id"`
	Status       string  `db:"status"`
	ErrorPhase   *string `db:"error_phase"`
	CPU          int     `db:"cpu"`
	MemoryMB     int     `db:"memory_mb"`
	ImageID      string  `db:"image_id"`
	SubnetID     string  `db:"subnet_id"`
	PrivateIP    string  `db:"private_ip"`
	RootVolumeID string  `db:"root_volume_id"`
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) compute.InstanceRepository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, inst *compute.Instance) error {
	db := sqlite.GetExt(ctx, r.db)

	model := instanceModel{
		ID:           string(inst.ID()),
		Name:         inst.Name(),
		OwnerID:      string(inst.OwnerID()),
		Status:       string(inst.Status()),
		CPU:          inst.CPU(),
		MemoryMB:     inst.MemoryMB(),
		ImageID:      string(inst.ImageID()),
		SubnetID:     string(inst.SubnetID()),
		PrivateIP:    inst.PrivateIP(),
		RootVolumeID: string(inst.RootVolumeID()),
	}

	if inst.ErrPhase() != nil {
		ep := string(*inst.ErrPhase())
		model.ErrorPhase = &ep
	}

	const query = `
INSERT INTO instances (
	id, name, owner_id, status, error_phase,
	cpu, memory_mb, image_id, subnet_id, private_ip, root_volume_id
) VALUES (
	:id, :name, :owner_id, :status, :error_phase, 
  :cpu, :memory_mb, :image_id, :subnet_id, :private_ip, :root_volume_id
) ON CONFLICT(id) DO UPDATE SET
	status = :status,
	error_phase = :error_phase,
	private_ip = :private_ip,
	root_volume_id = :root_volume_id
`
	_, err := sqlx.NamedExecContext(ctx, db, query, model)
	return err
}

func (r *repository) FindByID(ctx context.Context, id compute.InstanceID) (*compute.Instance, error) {
	const query = `SELECT * FROM instances WHERE id = ?`
	var m instanceModel
	if err := r.db.GetContext(ctx, &m, query, string(id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *repository) FindByOwnerID(ctx context.Context, ownerID user.UserID) ([]*compute.Instance, error) {
	const query = `SELECT * FROM instances WHERE owner_id = ?`
	var models []instanceModel
	if err := r.db.SelectContext(ctx, &models, query, string(ownerID)); err != nil {
		return nil, err
	}

	results := make([]*compute.Instance, len(models))
	for i, m := range models {
		results[i] = r.toEntity(&m)
	}
	return results, nil
}

func (r *repository) Delete(ctx context.Context, id compute.InstanceID) error {
	db := sqlite.GetExt(ctx, r.db)
	const query = `DELETE FROM instances WHERE id = ?`
	_, err := db.ExecContext(ctx, query, string(id))
	return err
}

func (r *repository) toEntity(m *instanceModel) *compute.Instance {
	// エンティティ側に Restore 用のファクトリ関数がある想定
	return compute.NewInstance(
		compute.InstanceID(m.ID),
		m.Name,
		user.UserID(m.OwnerID),
		compute.InstanceStatus(m.Status),
		(*compute.ErrPhase)(m.ErrorPhase),
		m.CPU,
		m.MemoryMB,
		image.ImageID(m.ImageID),
		network.SubnetID(m.SubnetID),
		m.PrivateIP,
		storage.VolumeID(m.RootVolumeID),
	)
}
