package data_test

import (
	"errors"
	"os"
	"reflect"
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
	validHash     = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA") // 64 bytes
	validTestFile = schema.File{
		Path:      "C:\\Music\\test.mp3",
		Hash:      validHash,
		MediaType: "mp3",
		Size:      1024,
		Mod:       time.Now(),
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

func withUnsetField[T any](obj *T, key string) T {
	if reflect.ValueOf(obj).IsNil() {
		panic("trying to unset a field on a nil element!")
	}
	if reflect.TypeOf(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		panic("trying to unset a field on an element that is not a structure!")
	}
	newObj := *obj
	reflect.ValueOf(&newObj).Elem().FieldByName(key).SetZero()
	return newObj
}

func TestSaveFileMissingFields(t *testing.T) {
	testCases := []struct {
		title       string
		file        schema.File
		expectedErr error
	}{
		{title: "Path", file: withUnsetField(&validTestFile, "Path"), expectedErr: data.ErrInvalidPath},
		{title: "Hash", file: withUnsetField(&validTestFile, "Hash"), expectedErr: data.ErrInvalidHash},
		{title: "MediaType", file: withUnsetField(&validTestFile, "MediaType"), expectedErr: data.ErrInvalidMediaType},
		{title: "Size", file: withUnsetField(&validTestFile, "Size"), expectedErr: data.ErrInvalidSize},
		{title: "Mod", file: withUnsetField(&validTestFile, "Mod"), expectedErr: data.ErrInvalidMod},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			err := data.SaveFile(db, tc.file)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
			}
		})
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

func TestSaveFileInvalidPath(t *testing.T) {
	invalidTestFileWindows := validTestFile
	invalidTestFileWindows.Path = "C:\\Windows\\Invalid:?Path"
	invalidTestFileUnix := validTestFile
	invalidTestFileUnix.Path = "/home/url/invalid:path"

	testCases := []struct {
		title string
		file  schema.File
	}{
		{title: "InvalidWindowsPath", file: invalidTestFileWindows},
		{title: "InvalidUnixPath", file: invalidTestFileUnix},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			err := data.SaveFile(db, tc.file)
			if !errors.Is(err, data.ErrInvalidPath) {
				t.Errorf("expected error %v, but got %v", data.ErrInvalidPath, err)
			}
		})
	}
}

func TestSaveFileInvalidHash(t *testing.T) {
	fileWithHashToSmall := validTestFile
	fileWithHashToSmall.Hash = make([]byte, 32)
	fileWithHashToBig := validTestFile
	fileWithHashToBig.Hash = make([]byte, 128)
	fileWithHashZero := validTestFile
	fileWithHashZero.Hash = make([]byte, 64) // all zeros

	testCases := []struct {
		title string
		file  schema.File
	}{
		{title: "HashToSmall", file: fileWithHashToSmall},
		{title: "HashToBig", file: fileWithHashToBig},
		{title: "HashZero", file: fileWithHashZero},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			err := data.SaveFile(db, tc.file)
			if !errors.Is(err, data.ErrInvalidHash) {
				t.Errorf("expected error %v, but got %v", data.ErrInvalidHash, err)
			}
		})
	}
}

