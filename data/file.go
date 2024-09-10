package data

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"

	"github.com/makl11/musiman/audio"
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

	if err := ValidatePath(file.Path); err != nil {
		return fmt.Errorf("%w: %w: \"%s\" is not a valid file path: %w", ErrInvalidPath, ErrInvalidArgumentValue, file.Path, err)
	}
	if len(file.Hash) != schema.HASH_SIZE {
		return fmt.Errorf("%w: %w: files content hash must consist of exactly %d bytes, but is %d bytes", ErrInvalidHash, ErrInvalidArgumentValue, schema.HASH_SIZE, len(file.Hash))
	}
	if isHashZero(file.Hash) {
		return fmt.Errorf("%w: %w: files content hash must not be zero", ErrInvalidHash, ErrInvalidArgumentValue)
	}
	if _, ok := audio.MUSIC_FILE_TYPES[file.MediaType]; !ok {
		return fmt.Errorf("%w: %w: unknown or unsupported media type: \"%s\"", ErrInvalidMediaType, ErrInvalidArgumentValue, file.MediaType)
	}

	return nil
}

// Based on: https://dwheeler.com/essays/fixing-unix-linux-filenames.html
func ValidatePath(path string) error {
	// Forbid ASCII control characters (bytes 1-31 and 127)
	for _, r := range path {
		if (r >= 0 && r <= 31) || r == 127 {
			return errors.New("path contains ASCII control characters")
		}
	}

	// Forbid leading "-"
	if len(path) > 0 && path[0] == '-' {
		return errors.New("path starts with a dash")
	}

	// Forbid filenames that aren’t a valid UTF-8 encoding
	if !utf8.ValidString(path) {
		return errors.New("path is not valid UTF-8")
	}

	// Forbid problematic characters
	if strings.ContainsAny(path, "*?[]\"<>|(){}&'!;") {
		return errors.New("path contains problematic characters")
	}

	// Forbid leading "~"
	if len(path) > 0 && path[0] == '~' {
		return errors.New("path starts with a tilde")
	}

	// Forbid "." and ".." as path elements
	if strings.Contains(path, "/../") || strings.Contains(path, "/./") || strings.Contains(path, "\\..\\") || strings.Contains(path, "\\.\\") {
		return errors.New("path contains \"..\" or \".\"")
	}

	// Forbid ":" except for after drive name in windows
	if regexp.MustCompile("^[A-Za-z]:").MatchString(path) {
		if strings.Contains(path[2:], ":") {
			return errors.New("path contains \":\" after drive name")
		}
	} else {
		if strings.Contains(path, ":") {
			return errors.New("path contains \":\"")
		}
	}

	return nil
}

func isHashZero(hash []byte) bool {
	for _, b := range hash {
		if b != 0 {
			return false
		}
	}
	return true
}
