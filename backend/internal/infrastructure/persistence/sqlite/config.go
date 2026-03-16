package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type DB interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func NewSqlxClient(cfg *Config) *sqlx.DB {
	// 接続
	db, err := sqlx.Open("sqlite", cfg.DBPath)
	if err != nil {
		panic(fmt.Errorf("failed to open sqlite: %w", err))
	}

	// 頼むから10コネクションでも整合性を保ってくれ
	// 厳しくなってきたら、Postgresへの移行も検討する
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// WALモードを有効化し、外部キー制約をONにする
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA synchronous=NORMAL;", // WALモード時はNORMALが推奨（性能と安全性のバランス）
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			panic(fmt.Sprintf("failed to set pragma '%s': %v", p, err))
		}
	}

	return db
}
