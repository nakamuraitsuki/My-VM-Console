package storage

import (
	"context"
	"database/sql"

	"example.com/m/internal/domain/storage"
	"github.com/jmoiron/sqlx"
)

type volumeModel struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	SizeGB int    `db:"size_gb"`
	Pool   string `db:"pool"`
	Owner  string `db:"owner"`
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) storage.Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(ctx context.Context, id storage.VolumeID) (*storage.Volume, error) {
	const query = `SELECT id, name, size_gb, pool, owner FROM volumes WHERE id = ?`
	var m volumeModel
	if err := r.db.GetContext(ctx, &m, query, string(id)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return storage.NewVolume(
		storage.VolumeID(m.ID),
		m.Name,
		m.SizeGB,
		m.Pool,
		m.Owner,
	), nil
}

func (r *repository) Save(ctx context.Context, volume *storage.Volume) error {
	model := volumeModel{
		ID:     string(volume.ID()),
		Name:   volume.Name(),
		SizeGB: volume.SizeGB(),
		Pool:   volume.Pool(),
		Owner:  volume.Owner(),
	}

	const query = `
INSERT INTO volumes (id, name, size_gb, pool, owner)
VALUES (:id, :name, :size_gb, :pool, :owner)
ON CONFLICT(id) DO UPDATE SET
	name = :name,
	size_gb = :size_gb,
	pool = :pool,
	owner = :owner
`

	_, err := r.db.NamedExecContext(ctx, query, model)
	return err
}

func (r *repository) Delete(ctx context.Context, id storage.VolumeID) error {
	const query = `DELETE FROM volumes WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, string(id))
	return err
}

