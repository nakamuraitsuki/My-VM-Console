package user

import (
	"context"

	"example.com/m/internal/domain/user"
	"example.com/m/internal/infrastructure/persistence/sqlite"
	"github.com/jmoiron/sqlx"
)

type userModel struct {
	ID string `db:"id"`
	// クォータ情報をフラットにマッピング
	MaxInstance int     `db:"quota_max_instance"`
	MaxCPU      int     `db:"quota_max_cpu"`
	MaxMemory   int     `db:"quota_max_memory"`
	Status      string  `db:"status"`
	ErrorPhase  *string `db:"error_phase"`
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) user.UserRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(ctx context.Context, id user.UserID) (*user.UserPersistentData, error) {
	const query = `
SELECT id, quota_max_instance, quota_max_cpu, quota_max_memory, status, error_phase
FROM users
WHERE id = ?
`
	var um userModel
	// sqlx.GetContext で一気にマッピング
	if err := r.db.GetContext(ctx, &um, query, string(id)); err != nil {
		return nil, err
	}

	return &user.UserPersistentData{
		ID: user.UserID(um.ID),
		Quota: user.UsageQuota{
			MaxInstance: um.MaxInstance,
			MaxCPU:      um.MaxCPU,
			MaxMemory:   um.MaxMemory,
		},
		Status:     user.UserStatus(um.Status),
		ErrorPhase: (*user.FailedPhase)(um.ErrorPhase), // シンプルなキャスト
	}, nil
}

func (r *repository) Save(ctx context.Context, usr *user.User) error {
	db := sqlite.GetExt(ctx, r.db)
	// 永続化モデルへの詰め替え
	model := userModel{
		ID:          string(usr.ID()),
		MaxInstance: usr.Quota().MaxInstance,
		MaxCPU:      usr.Quota().MaxCPU,
		MaxMemory:   usr.Quota().MaxMemory,
		Status:      string(usr.Status()),
	}
	if usr.ErrPhase() != nil {
		p := string(*usr.ErrPhase())
		model.ErrorPhase = &p
	}

	const query = `
INSERT INTO users (id, quota_max_instance, quota_max_cpu, quota_max_memory, status, error_phase)
VALUES (:id, :quota_max_instance, :quota_max_cpu, :quota_max_memory, :status, :error_phase)
ON CONFLICT(id) DO UPDATE SET
	quota_max_instance = :quota_max_instance,
	quota_max_cpu = :quota_max_cpu,
	quota_max_memory = :quota_max_memory,
	status = :status,
	error_phase = :error_phase
`
	// NamedExecContext でモデルをそのまま渡す
	_, err := sqlx.NamedExecContext(ctx, db, query, model)
	return err
}

func (r *repository) Delete(ctx context.Context, id user.UserID) error {
	db := sqlite.GetExt(ctx, r.db)
	const query = `DELETE FROM users WHERE id = ?`
	_, err := db.ExecContext(ctx, query, string(id))
	return err
}
