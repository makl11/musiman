package data

import (
	"github.com/jmoiron/sqlx"

	"github.com/makl11/musiman/data/schema"
)

func SaveFile(db *sqlx.DB, file schema.File) error {
	_, err := db.NamedExec(`INSERT INTO files (path, hash, media_type, size, mod) VALUES (:path, :hash, :media_type, :size, :mod)`, file)
	return err
}
