package data_test

import (
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"github.com/makl11/musiman/data"
	"github.com/makl11/musiman/data/schema"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	goose.SetLogger(goose.NopLogger())
	fs := os.DirFS(".")
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("sqlite"); err != nil {
		t.Fatalf("failed to set sqlite dialect for goose: %v", err)
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		t.Fatalf("failed to apply goose migrations: %v", err)
	}

	return db
}

var (
	validTestFile = schema.File{
		Path:      "test.mp3",
		Hash:      make([]byte, 64),
		MediaType: "mp3",
		Size:      1024,
		Mod:       time.Now().Format("ISO8601"),
	}
)

func TestSaveFileSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := data.SaveFile(db, validTestFile)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM files WHERE path = ?", validTestFile.Path)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 record, but got %d", count)
	}
}

func TestSaveFileDuplicatePath(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := data.SaveFile(db, validTestFile)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	err = data.SaveFile(db, validTestFile)
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}
