package data

import (
	"embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "data/data.db")
	if err != nil {
		return nil, err
	}

	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return nil, err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return nil, err
	}

	return db, nil
}
