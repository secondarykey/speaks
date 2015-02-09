package db

import (
	"database/sql"
)

type Role struct {
	Id   int
	Key  string
	Name string
}

func createRoleTable() error {
	_, err := Exec("DROP TABLE if exists Role")
	if err != nil {
		return err
	}
	_, err = Exec("CREATE TABLE Role(key text PRIMARY KEY,name text)")
	return err
}

func insertRole(tx *sql.Tx, name string, key string) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into Role(name,key) values(?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(name, key)
}
