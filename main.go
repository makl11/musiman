package main

import (
	"os"
	"path/filepath"

	"github.com/makl11/musiman/scanner"
)

func main() {
	var err error
	scanRoot, _ := filepath.Abs(".")
	if len(os.Args) > 1 {
		scanRoot, err = filepath.Abs(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	scanner.ScanDirForMusic(scanRoot)
}
