package scanner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/liamg/magic"
)

var MUSIC_FILE_TYPES = map[string]bool{
	// https://en.wikipedia.org/wiki/MP3
	"mp3": true,
	// https://en.wikipedia.org/wiki/Ogg
	"ogg": true,
	"oga": true,
	// https://en.wikipedia.org/wiki/Windows_Media_Audio
	"wma": true,
	// https://en.wikipedia.org/wiki/Free_Lossless_Audio_Codec
	"flac": true,
	// https://en.wikipedia.org/wiki/Waveform_Audio_File_Format
	"wav": true,
	// https://en.wikipedia.org/wiki/Audio_Interchange_File_Format
	"aiff": true,
	"aif":  true,
	"aifc": true,
	"snd":  true,
	"iff":  true,
}

// TODO: skip path if it matches any of the ignorePaths
func ScanDirForMusic(scanRoot string, minSize uint64, ignorePaths []string) error {
	fileSystem := os.DirFS(scanRoot)
	var buf []byte = make([]byte, 1024) // size recommended here: https://pkg.go.dev/github.com/liamg/magic#Lookup
	return fs.WalkDir(fileSystem, ".", func(path string, entry fs.DirEntry, err error) error {
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

			if MUSIC_FILE_TYPES[fileType.Extension] {
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
