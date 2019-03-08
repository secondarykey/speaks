package db

import (
	"database/sql"
	"log"
)

type UserRole struct {
	UserId  int
	RoleKey string
}

func createUserRoleTable() error {
	_, err := Exec("CREATE TABLE UserRole(user_id int,role_key text)")
	return err
}

func dropUserRoleTable() error {
	_, err := Exec("DROP TABLE if exists UserRole")
	return err
}

func InsertUserRole(tx *sql.Tx, userId int, roleKey string) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into UserRole(user_id,role_key) values(?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(userId, roleKey)
}

func SelectAdminUsers() ([]UserRole, error) {
	sql := "select user_id from UserRole where role_key = ?"
	rows, err := inst.Query(sql, RoleAdmin)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	roles := make([]UserRole, 0)
	for rows.Next() {
		role := UserRole{}
		rows.Scan(&role.UserId)
		roles = append(roles, role)
	}
	return roles, nil
}

func SelectUserRole(userId int) ([]UserRole, error) {
	sql := "select role_key from UserRole where user_id = ?"
	rows, err := inst.Query(sql, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	roles := make([]UserRole, 0)
	for rows.Next() {
		role := UserRole{}
		rows.Scan(&role.RoleKey)
		roles = append(roles, role)
	}
	return roles, nil
}

func InitUserRole(tx *sql.Tx, userId int) error {

	_, err := InsertUserRole(tx, userId, RoleAdmin)
	if err != nil {
		return err
	}

	_, err = InsertUserRole(tx, userId, RoleSpeaker)
	if err != nil {
		return err
	}

	return nil
}
