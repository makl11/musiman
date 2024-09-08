package data

import (
	"context"
	"embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/makl11/musiman/context_keys"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitDb(cmd *cobra.Command, args []string) error {
	db, err := sqlx.Connect("sqlite3", "data/data.db")
	if err != nil {
		return err
	}

	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), context_keys.DB, db))
	return nil
}
