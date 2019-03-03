package db

import (
	"database/sql"
)

type Role struct {
	Id   int
	Key  string
	Name string
}

const (
	RoleAdmin   = "Administrator"
	RoleSpeaker = "Speaker"
)

func NewRoleMap() RoleMap {
	return map[string]bool{
		RoleAdmin:   false,
		RoleSpeaker: false,
	}
}

func createRoleTable() error {
	_, err := Exec("CREATE TABLE Role(key text PRIMARY KEY,name text)")
	return err
}

func dropRoleTable() error {
	_, err := Exec("DROP TABLE if exists Role")
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

func InitRole(tx *sql.Tx) error {
	_, err := insertRole(tx, "Administrator", RoleAdmin)
	if err != nil {
		return err
	}
	_, err = insertRole(tx, "Speaker", RoleSpeaker)
	if err != nil {
		return err
	}
	return nil
}
