package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func createUserTable() error {
	_, err := Exec("DROP TABLE if exists User")
	if err != nil {
		return err
	}
	_, err = Exec("CREATE TABLE User(id INTEGER PRIMARY KEY AUTOINCREMENT,name text,email text,password text)")
	return err
}

func insertUser(tx *sql.Tx, name string, email string, password string) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into User(name,email,password) values(?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	pswd := CreateMD5(password)
	return stmt.Exec(name, email, pswd)
}

func CreateMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func UpdateUser(u *User) error {
	_, err := inst.Exec("update user set name=?,email=?,password=? where id = ?",
		u.Name, u.Email, u.Password, u.Id)
	return err
}

func SelectUser(email, pswd string) (*User, error) {

	user := &User{}
	err := inst.QueryRow("select id, name ,email,password from user where email = ?", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	if user.Password == CreateMD5(pswd) {
		return user, nil
	}
	return nil, errors.New("パスワードが違うよ")
}
