package db

import (
	"database/sql"
)

type Category struct {
	Id          int
	Key         string
	Name        string
	Description string
}

func createCategoryTable() error {
	_, err := Exec("DROP TABLE if exists Category")
	if err != nil {
		return err
	}
	_, err = Exec("CREATE TABLE Category(id INTEGER PRIMARY KEY AUTOINCREMENT,key text,name text,description text)")
	return err
}

func InsertCategory(key, name, desc string) (sql.Result, error) {
	return inst.Exec("insert into Category(key,name,desc) values(?, ?, ?)", key, name, desc)
}
