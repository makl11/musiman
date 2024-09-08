package main

import (
	"github.com/makl11/musiman/cmd"
	"github.com/makl11/musiman/data"
)

func main() {
	db, err := data.InitDb()
	if err != nil {
		panic(err)
	}
	db.Close()
	cmd.Execute()
}
