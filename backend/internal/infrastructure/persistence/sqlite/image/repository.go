package image

import (
	"context"
	"database/sql"

	"example.com/m/internal/domain/image"
	"github.com/jmoiron/sqlx"
)

type imageModel struct {
	ID          string `db:"id"`
	Alias       string `db:"alias"`
	Fingerprint string `db:"fingerprint"`
	ServerURL   string `db:"server_url"`
	Protocol    string `db:"protocol"`
	IsPublic    bool   `db:"is_public"`
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) image.Repository {
	return &repository{db: db}
}

func (r *repository	) FindByAlias(ctx context.Context, alias string) (*image.Image, error) {
	const query = `SELECT * FROM images WHERE alias = ?`
	var m imageModel
	if err := r.db.GetContext(ctx, &m, query, alias); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *repository) FindAll(ctx context.Context) ([]*image.Image, error) {
	const query = `SELECT * FROM images`
	var models []imageModel
	if err := r.db.SelectContext(ctx, &models, query); err != nil {
		return nil, err
	}

	results := make([]*image.Image, len(models))
	for i, m := range models {
		results[i] = r.toEntity(&m)
	}
	return results, nil
}

func (r *repository) Save(ctx context.Context, img *image.Image) error {
	model := imageModel{
		ID:          string(img.ID()),
		Alias:       img.Alias(),
		Fingerprint: img.Fingerprint(),
		ServerURL:   img.ServerURL(),
		Protocol:    img.Protocol(),
		IsPublic:    img.IsPublic(),
	}

	const query = `
INSERT INTO images (id, alias, fingerprint, server_url, protocol, is_public)
VALUES (:id, :alias, :fingerprint, :server_url, :protocol, :is_public)
ON CONFLICT(id) DO UPDATE SET
	alias = :alias,
	fingerprint = :fingerprint,
	server_url = :server_url,
	protocol = :protocol,
	is_public = :is_public
`
	_, err := r.db.NamedExecContext(ctx, query, model)
	return err
}

func (r *repository) toEntity(m *imageModel) *image.Image {
	return image.NewImage(
		image.ImageID(m.ID),
		m.Alias,
		m.Fingerprint,
		m.ServerURL,
		m.Protocol,
		m.IsPublic,
	)
}
