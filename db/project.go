package db

import (
	"database/sql"
)

type Project struct {
	Key         string
	Name        string
	Seq         int
	Description string
}

const (
	DefaultProject = "Speaks"
)

func (p Project) TableName() string {
	return "Project"
}

func (p Project) Create() error {
	_, err := Exec("CREATE TABLE Project(key text PRIMARY KEY,name text,seq int,description text)")
	return err
}

func (p Project) Drop() error {
	_, err := Exec("DROP TABLE if exists Role")
	return err
}

func (p Project) Init(tx *sql.Tx) error {

	p.Key = DefaultProject
	p.Name = "Speaks"
	p.Seq = 0
	p.Description = "Everyone's Project"
	_, err := p.Insert(tx)
	return err
}

func (p Project) Insert(tx *sql.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into Project(key,name,seq,description) values(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(p.Key, p.Name, p.Seq, p.Description)
}

func createProjectTable() error {
	p := Project{}
	return p.Create()
}

func dropProjectTable() error {
	p := Project{}
	return p.Drop()
}

func InitProject(tx *sql.Tx) error {
	p := Project{}
	return p.Init(tx)
}
