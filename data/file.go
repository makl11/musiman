package data

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/makl11/musiman/data/schema"
)

func SaveFile(db *sqlx.DB, file schema.File) error {
	if err := ValidateFile(file); err != nil {
		return err
	}

	_, err := db.NamedExec(`INSERT INTO files (path, hash, media_type, size, mod) VALUES (:path, :hash, :media_type, :size, :mod)`, file)
	return err
}

var (
	ErrMissingArgumentValue = errors.New("missing argument value")
	ErrInvalidArgumentValue = errors.New("invalid argument value")

	ErrInvalidPath      = errors.New("invalid path")
	ErrInvalidHash      = errors.New("invalid hash")
	ErrInvalidMediaType = errors.New("invalid media type")
	ErrInvalidSize      = errors.New("invalid size")
	ErrInvalidMod       = errors.New("invalid mod")
)

func ValidateFile(file schema.File) error {
	if file.Path == "" {
		return fmt.Errorf("%w: %w: path must not be empty", ErrInvalidPath, ErrMissingArgumentValue)
	}
	if len(file.Hash) == 0 {
		return fmt.Errorf("%w: %w: hash must not be empty", ErrInvalidHash, ErrMissingArgumentValue)
	}
	if file.MediaType == "" {
		return fmt.Errorf("%w: %w: media type must not be empty", ErrInvalidMediaType, ErrMissingArgumentValue)
	}
	if file.Size == 0 {
		return fmt.Errorf("%w: %w: size must not be zero", ErrInvalidSize, ErrMissingArgumentValue)
	}
	if file.Mod.IsZero() {
		return fmt.Errorf("%w: %w: mod time must not be zero", ErrInvalidMod, ErrMissingArgumentValue)
	}
	return nil
}
