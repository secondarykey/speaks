package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

var inst *sql.DB
var tx *sql.Tx

type schemaError struct {
	code    int
	message string
}

func (s *schemaError) Error() string {
	return fmt.Sprintf("%d:%s", s.code, s.message)
}

func NewSchemaError(code int, msg string) *schemaError {
	return &schemaError{
		code:    code,
		message: msg,
	}
}

const schemaVersion = "0.1"

func check(path, ver string) (string, *schemaError) {

	//%sがあるか？
	pArr := strings.Split(path, "%s")
	if len(pArr) != 2 {
		return "", NewSchemaError(-1, "Error:database path is '%s' requid["+path+"]")
	}

	rpath := fmt.Sprintf(path, schemaVersion)
	//存在するか？
	_, err := os.Stat(rpath)
	//versionが一緒か？
	if ver == schemaVersion || ver == "test" {
		if err == nil {
			return rpath, nil
		}
		return rpath, NewSchemaError(0, "Create database")
	}

	if err == nil {
		return rpath, nil
	}

	//code 0 ってなんだ？

	return rpath, NewSchemaError(0, "Warning:Program version,TOML file version")
}

func Listen(path, version string) error {

	var err error
	rp, scErr := check(path, version)

	cFlag := true
	if scErr != nil {
		if scErr.code == 0 {
			cFlag = false
		} else {
			return scErr
		}
	}

	inst, err = sql.Open("sqlite3", rp)
	if err != nil {
		return err
	}

	if cFlag {
		return nil
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
	err = createCategoryTable()
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
