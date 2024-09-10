package schema

import "time"

const HASH_SIZE = 64

type File struct {
	Path      string
	Hash      []byte // (schema.HASH_SIZE bytes) must be unsized for storage driver compatibility
	MediaType string `db:"media_type"`
	Size      uint
	Mod       time.Time
}
