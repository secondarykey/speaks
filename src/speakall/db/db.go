package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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

func createInitTable() error {
	err := createUserTable()
	if err != nil {
		return nil
	}
	err = createRoleTable()
	if err != nil {
		return nil
	}
	return createUserRoleTable()
}

func insertInitTable() error {
	tx, err := inst.Begin()
	if err != nil {
		return err
	}

	rslt, err := insertUser(tx, "SpeakAll管理者", "admin@localhost", "password")
	if err != nil {
		return err
	}
	userId, _ := rslt.LastInsertId()

	rslt, err = insertRole(tx, "管理者", "Admin")
	rslt, err = insertRole(tx, "議題編集者", "Chairman")

	rslt, err = insertRole(tx, "発言者", "Speaker")
	rslt, err = insertRole(tx, "閲覧者", "Viewer")
	if err != nil {
		return err
	}

	rslt, err = insertUserRole(tx, int(userId), "Admin")
	rslt, err = insertUserRole(tx, int(userId), "Chairman")
	rslt, err = insertUserRole(tx, int(userId), "Speaker")
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func Exec(sql string) (sql.Result, error) {
	return inst.Exec(sql)
}

/*
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
*/
