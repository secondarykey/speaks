package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var inst *sql.DB

func Listen(path string) error {
	var err error
	inst, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	err = createInitTable()
	if err != nil {
		return err
	}

	return insertInitTable()
}

func Exec(sql string) (sql.Result, error) {
	return inst.Exec(sql)
}

func createInitTable() {
	err := createUserTable()
	log.Println(err)
	err = createRoleTable()
	log.Println(err)
	err = createUserRoleTable()
	log.Println(err)
}

func insertInitTable() error {
	tx, err := inst.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	rslt, err := insertUser(tx, "SpeakAll管理者", "admin@local.host", "password")
	if err != nil {
		return err
	}
	userId, _ := rslt.LastInsertId()

	rslt, err = insertRole(tx, "管理者", "Admin")
	if err != nil {
		return err
	}

	rslt, err = insertUserRole(tx, int(userId), "Admin")
	if err != nil {
		return err
	}

	tx.Commit()
}

func Select() {

	rows, err := inst.Query("select id, name from foo")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		println(id, name)
	}
}
