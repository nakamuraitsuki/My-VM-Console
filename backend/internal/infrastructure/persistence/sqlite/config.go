package sqlite

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqlxClient(cfg *Config) *sqlx.DB {
	// 接続
	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		panic(fmt.Errorf("failed to open sqlite: %w", err))
	}

	// 接続設定（SQLiteの定石）
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
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