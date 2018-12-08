package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/prest/adapters"
)

func Loaddata(filename string, db adapters.Adapter) {
	filedata, _ := ioutil.ReadFile(filename)
	registers := []Register{}
	err := json.Unmarshal(filedata, &registers)
	if err != nil {
		log.Println(err)
		return
	}
	for _, register := range registers {
		names, placeholders, values, err := ParseInsertRegister(register)
		if err != nil {
			log.Println(err)
			return
		}
		sql := db.InsertSQL(register.Database, register.Schema, register.Table, names, placeholders)
		//fmt.Println(sql, values)
		sc := db.Insert(sql, values...)
		if sc.Err() != nil {
			log.Println(sc.Err())
			return
		}
	}
	fmt.Println(filename, " loaded")
}
