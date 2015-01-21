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

	createInitTable()
	insertInitTable()

	return err
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

func insertInitTable() {
	tx, err := inst.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	rslt, err := insertUser(tx, "SpeakAll管理者", "admin@local.host", "password")
	userId, _ := rslt.LastInsertId()

	rslt, err = insertRole(tx, "管理者", "Admin")
	log.Println(err)

	rslt, err = insertUserRole(tx, int(userId), "Admin")
	log.Println(err)

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
