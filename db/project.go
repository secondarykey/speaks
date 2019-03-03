package db

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
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
	_, err := Exec("DROP TABLE if exists Project")
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

func SelectProjects() ([]Project, error) {

	rows, err := inst.Query("select key,name,description from Project order by seq asc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]Project, 0)
	for rows.Next() {
		p := Project{}
		rows.Scan(&p.Key, &p.Name, &p.Description)
		projects = append(projects, p)
	}
	return projects, nil
}

func InsertProject(name, desc string) error {

	tx, err := inst.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	p := Project{}
	p.Key = uuid.NewV4().String()
	p.Name = name
	p.Seq = 1
	p.Description = desc

	_, err = p.Insert(tx)
	if err != nil {
		return err
	}

	err = InsertDefaultMember(tx, p.Key)
	if err != nil {
		return err
	}

	c := Category{}
	err = c.InsertDefaultCategory(tx, p.Key)
	if err != nil {
		return err
	}

	return tx.Commit()
}
