package schema

type File struct {
	Path      string
	Hash      []byte
	MediaType string `db:"media_type"`
	Size      uint
	Mod       string
}
