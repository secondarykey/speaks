package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
	. "github.com/secondarykey/speaks/config"
	"github.com/secondarykey/speaks/logic"
)

type RoleMap map[string]bool
type ProjectRoleMap map[string]RoleMap

type User struct {
	Id       int
	Name     string
	Email    string
	Password string

	CurrentProject Project
	Roles          RoleMap
	Projects       []Project
	ProjectRoles   ProjectRoleMap
}

func createUserTable() error {
	_, err := Exec("CREATE TABLE User(id INTEGER PRIMARY KEY AUTOINCREMENT,name text,email text,password text)")
	return err
}

func dropUserTable() error {
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
	if err != nil {
		return err
	}
	userId, _ := rslt.LastInsertId()

	_, err = InsertUserRole(tx, int(userId), DefaultProject)
	if err != nil {
		return err
	}

	err = InsertDefaultMember(tx, int(userId))
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

func InsertLDAPUser(result *logic.LDAPResult) error {

	tx, err := inst.Begin()
	if err != nil {
		return err
	}

	//User
	rslt, err := InsertUser(tx, result.Name, result.LoginName, "")
	if err != nil {
		return err
	}
	userId, _ := rslt.LastInsertId()

	//Role
	_, err = InsertUserRole(tx, int(userId), RoleSpeaker)
	if err != nil {
		return err
	}

	//Speaks Project Member
	err = InsertDefaultMember(tx, int(userId))
	if err != nil {
		return err
	}

	return tx.Commit()
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
		log.Println(err)
		return nil, err
	}

	if Config.LDAP.Use {
		return user, nil
	}

	if user.Password != CreateMD5(pswd) {
		log.Println(err)
		return nil, fmt.Errorf("Auth Error")
	}

	return user, nil
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

func (u *User) IsLogin() bool {
	if u == nil || u.Name == "" {
		return false
	}
	return true
}

func (u *User) IsAdmin() bool {
	if u.IsLogin() {
		return u.Roles[RoleAdmin]
	}
	return false
}

func (u *User) IsSpeaker() bool {
	if u.IsLogin() {
		return u.Roles[RoleSpeaker]
	}
	return false
}

func (u *User) See(projectKey string) bool {
	if u.IsLogin() {
		_, ok := u.ProjectRoles[projectKey]
		return ok
	}
	return false
}

func (u *User) IsManager() bool {
	if u.IsLogin() {
		roles, ok := u.ProjectRoles[u.CurrentProject.Key]
		if !ok {
			return false
		}
		return roles[MemberManager]
	}
	return false
}

func (u *User) IsEditor() bool {
	if u.IsLogin() {
		roles, ok := u.ProjectRoles[u.CurrentProject.Key]
		if !ok {
			return false
		}
		return roles[MemberEditor]
	}
	return false
}

func (u *User) IsViewer() bool {
	if u.IsLogin() {
		roles, ok := u.ProjectRoles[u.CurrentProject.Key]
		if !ok {
			return false
		}
		return roles[MemberViewer]
	}
	return false
}

func (u *User) Init(tx *sql.Tx) error {
	pwd := CreateMD5("p@ssword")
	rslt, err := InsertUser(tx, "Speaks Administrator", "admin@localhost", pwd)
	if err != nil {
		return err
	}
	userId, err := rslt.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = int(userId)
	return nil
}

func InitUser(tx *sql.Tx) (int, error) {
	u := User{}
	err := u.Init(tx)
	if err != nil {
		return -1, err
	}
	return u.Id, nil
}
