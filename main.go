package main

import (
	"os"

	"fmt"

	"github.com/joelmdesouza/prest-dumpdata/database"
	"github.com/prest/adapters/postgres"
	"github.com/prest/config"
)

func main() {
	config.Load()
	postgres.Load()
	db := config.PrestConf.Adapter

	action := "dumpdata"
	filename := "fixture.json"

	args := os.Args[1:]

	if len(args) > 0 {
		action = args[0]
	}

	if len(args) > 1 {
		filename = args[1]
	}

	switch action {
	case "dumpdata":
		database.Dumpdata(filename, db)
	case "loaddata":
		database.Loaddata(filename, db)
	default:
		fmt.Println("Invalid Action")
	}
}
