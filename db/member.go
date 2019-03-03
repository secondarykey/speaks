package db

import (
	"database/sql"
)

type Member struct {
	Project string
	UserId  int
	Role    string
}

const (
	MemberManager = "Manager"
	MemberEditor  = "Editor"
	MemberSpeaker = "Speaker"
)

func NewMemberRole() RoleMap {
	return map[string]bool{
		MemberManager: false,
		MemberEditor:  false,
		MemberSpeaker: false,
	}
}

func (p Member) TableName() string {
	return "Member"
}

func (p Member) Create() error {
	_, err := Exec("CREATE TABLE Member(project text,user_id int,role text)")
	return err
}

func (p Member) Drop() error {
	_, err := Exec("DROP TABLE if exists Member")
	return err
}

func (m *Member) Init(tx *sql.Tx) error {

	m.Role = MemberManager
	_, err := m.Insert(tx)
	if err != nil {
		return err
	}

	m.Role = MemberEditor
	_, err = m.Insert(tx)
	if err != nil {
		return err
	}
	m.Role = MemberSpeaker
	_, err = m.Insert(tx)
	if err != nil {
		return err
	}

	return nil
}

func (p Member) Insert(tx *sql.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into Member(project,user_id,role) values(?, ?,?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(p.Project, p.UserId, p.Role)
}

func InsertDefaultMember(tx *sql.Tx, project string) error {

	uroles, err := SelectAdminUsers()
	if err != nil {
		return err
	}

	m := Member{}
	m.Project = project
	m.Role = MemberManager
	for _, ur := range uroles {
		m.UserId = ur.UserId
		_, err = m.Insert(tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func createMemberTable() error {
	p := Member{}
	return p.Create()
}

func dropMemberTable() error {
	p := Member{}
	return p.Drop()
}

func InitMember(tx *sql.Tx, id int) error {
	p := Member{}
	p.Project = DefaultProject
	p.UserId = id
	return p.Init(tx)
}

func SelectMember(id int) ([]Member, error) {

	rows, err := inst.Query("select project,user_id,role from Member Where user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]Member, 0)
	for rows.Next() {
		mem := Member{}
		rows.Scan(&mem.Project, &mem.UserId, &mem.Role)
		members = append(members, mem)
	}
	return members, nil
}
