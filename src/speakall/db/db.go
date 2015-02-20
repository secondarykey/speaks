package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

var inst *sql.DB
var tx *sql.Tx

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
		return err
	}
	err = createRoleTable()
	if err != nil {
		return err
	}
	err = createUserRoleTable()
	if err != nil {
		return err
	}
	err = createMessageTable()
	if err != nil {
		return err
	}
	return nil
}

func begin() error {
	var err error
	tx, err = inst.Begin()
	if err != nil {
		return err
	}
	return nil
}

func rollback() error {
	if tx != nil {
		err := tx.Rollback()
		tx = nil
		return err
	}
	return errors.New("not use Tx")
}

func commit() error {
	err := tx.Commit()
	tx = nil
	return err
}

func insertInitTable() error {

	err := begin()
	defer rollback()
	if err != nil {
		return err
	}

	rslt, err := insertUser(tx, "SpeakAll管理者", "admin@localhost", "password")
	if err != nil {
		return err
	}
	userId, _ := rslt.LastInsertId()

	rslt, err = insertRole(tx, "管理者", "Admin")
	if err != nil {
		return err
	}
	rslt, err = insertRole(tx, "議題編集者", "Chairman")
	if err != nil {
		return err
	}

	rslt, err = insertRole(tx, "発言者", "Speaker")
	if err != nil {
		return err
	}
	rslt, err = insertRole(tx, "閲覧者", "Viewer")
	if err != nil {
		return err
	}

	rslt, err = insertUserRole(tx, int(userId), "Admin")
	if err != nil {
		return err
	}

	rslt, err = insertUserRole(tx, int(userId), "Chairman")
	if err != nil {
		return err
	}

	rslt, err = insertUserRole(tx, int(userId), "Speaker")
	if err != nil {
		return err
	}

	commit()
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
