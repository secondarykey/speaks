package db

import (
	"database/sql"
)

type UserRole struct {
	UserId  int
	RoleKey string
}

func createUserRoleTable() error {
	_, err := Exec("DROP TABLE if exists UserRole")
	if err != nil {
		return err
	}
	_, err = Exec("CREATE TABLE UserRole(user_id int,role_key text)")
	return err
}

func InsertUserRole(tx *sql.Tx, userId int, roleKey string) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into UserRole(user_id,role_key) values(?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rslt, err := stmt.Exec(userId, roleKey)
	return rslt, err
}
