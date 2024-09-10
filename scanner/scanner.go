package scanner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/liamg/magic"

	"github.com/makl11/musiman/audio"
)

// TODO: skip path if it matches any of the ignorePaths
func ScanDirForMusic(scanRoot string, minSize uint64, ignorePaths []string) error {
	fileSystem := os.DirFS(scanRoot)
	var buf []byte = make([]byte, 1024) // size recommended here: https://pkg.go.dev/github.com/liamg/magic#Lookup
	return fs.WalkDir(fileSystem, ".", func(path string, entry fs.DirEntry, err error) error {
		// check if path is in ignorePaths
		for _, ignorePath := range ignorePaths {
			if strings.HasPrefix(filepath.Clean(path), filepath.Clean(ignorePath)) {
				return nil
			}
		}

		if !entry.IsDir() {
			fileInfo, err := fs.Stat(fileSystem, path)
			if err != nil {
				return nil
			}

			if fileInfo.Size() < int64(minSize) {
				return nil
			}

			f, err := fileSystem.Open(path)
			if err != nil {
				return nil
			}
			bytesRead, err := f.Read(buf)
			f.Close()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}

			if bytesRead == 0 {
				return nil
			}
			fileType, err := magic.LookupSync(buf[:bytesRead])
			if err != nil {
				if err == magic.ErrUnknown {
					return nil
				}
				return err
			}

			if fileType == nil {
				panic("filetype is nil")
			}

			if audio.MUSIC_FILE_TYPES[fileType.Extension] {
				fmt.Print(filepath.Join(scanRoot, path), "\t\t")
				if fileType.MIME != "" {
					fmt.Printf("%s\t", fileType.MIME)
				}
				fmt.Printf("File type: %s\t", fileType.Description)
				fmt.Println()
			}

		}
		return nil
	})
}
