package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/prest/adapters"
)

func Dumpdata(filename string, db adapters.Adapter) {
	sc := db.Query(TablesSelect)
	err := sc.Err()
	if err != nil {
		log.Println(err)
		return
	}

	tables := []Table{}
	_, err = sc.Scan(&tables)
	if err != nil {
		log.Println(err)
		return
	}

	registers := []Register{}
	for _, table := range tables {
		sql := fmt.Sprintf(`SELECT * FROM "%s"."%s"."%s"`, table.Database, table.Schema, table.Name)
		sc := db.Query(sql)
		err = sc.Err()
		if err != nil {
			log.Println(err)
			return
		}

		var rows []map[string]interface{}
		_, err = sc.Scan(&rows)
		if err != nil {
			log.Println(err)
			return
		}

		for _, row := range rows {
			register := Register{}
			register.Database = "prest"
			register.Schema = table.Schema
			register.Table = table.Name
			register.Fields = row
			registers = append(registers, register)
		}
	}

	//json, err := json.Marshal(registers)
	json, err := json.MarshalIndent(registers, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	err = ioutil.WriteFile(filename, json, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(filename, " generated")
}
