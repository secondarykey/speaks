package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	uuid "github.com/satori/go.uuid"
)

type RoleMap map[string]bool
type User struct {
	Id       int
	Name     string
	Email    string
	Password string
	Roles    RoleMap
}

func createUserTable() error {
	_, err := Exec("CREATE TABLE User(id INTEGER PRIMARY KEY AUTOINCREMENT,name text,email text,password text)")
	return err
}

func deleteUserTable() error {
	_, err := Exec("DROP TABLE if exists User")
	return err
}

func InsertUser(tx *sql.Tx, name string, email string, password string) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into User(name,email,password) values(?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(name, email, password)
}

func CreateUser(u *User) error {
	tx, err := inst.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	rslt, err := InsertUser(tx, u.Name, u.Email, uuid.NewV4().String())
	userId, _ := rslt.LastInsertId()
	rslt, err = InsertUserRole(tx, int(userId), "Speaker")
	if err != nil {
		return err
	}
	return tx.Commit()
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
		roles, err := selectUserRole(user.Id)
		if err != nil {
			return nil, err
		}
		user.Roles = map[string]bool{
			"Admin":    false,
			"Chairman": false,
			"Speaker":  false,
		}
		for _, elm := range roles {
			user.Roles[elm.RoleKey] = true
		}
		return user, nil
	}

	return nil, errors.New("パスワードが違うよ")
}

func SelectPassword(pswd string) (*User, error) {
	user := &User{}
	err := inst.QueryRow("select id, name ,email,password from user where password = ?", pswd).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserName(id int) (string, error) {
	var name string
	err := inst.QueryRow("select name from user where id = ?", id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func SelectAllUser() ([]*User, error) {
	sql := "select id,name,email,password from user"
	rows, err := inst.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		users = append(users, user)
	}
	return users, nil
}

func (u *User) IsAdmin() bool {
	return u.Roles["Admin"]
}

func (u *User) IsChairman() bool {
	return u.Roles["Chairman"]
}

func (u *User) IsSpeaker() bool {
	return u.Roles["Speaker"]
}
