package db

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
)

const (
	DefaultCategory = "Dashboard"
)

type Category struct {
	Id          int
	Key         string
	Project     string
	Name        string
	Description string
}

func (c Category) Create() error {
	_, err := Exec("CREATE TABLE Category(id INTEGER PRIMARY KEY AUTOINCREMENT,key text,name text,project text,description text)")
	return err
}

func (c Category) Drop() error {
	_, err := Exec("DROP TABLE if exists Category")
	return err
}

func createCategoryTable() error {
	c := Category{}
	return c.Create()
}

func dropCategoryTable() error {
	c := Category{}
	return c.Drop()
}

func (c *Category) Init(tx *sql.Tx) error {
	return c.InsertDefaultCategory(tx, DefaultProject)
}

func (c *Category) InsertDefaultCategory(tx *sql.Tx, project string) error {

	c.Key = DefaultCategory
	c.Name = "Dashboard"
	c.Project = project
	c.Description = "First Category"

	_, err := c.Insert(tx)
	if err != nil {
		return err
	}
	return nil
}

func (c *Category) Insert(tx *sql.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare("insert into Category(key,name,project,description) values(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(c.Key, c.Name, c.Project, c.Description)
}

func InitCategory(tx *sql.Tx) error {
	c := Category{}
	return c.Init(tx)
}

func GenerateCategoryKey() (string, error) {
	for {
		genKey := uuid.NewV4().String()
		var key string
		err := inst.QueryRow("select key from Category where key = ?", genKey).Scan(key)
		switch {
		case err == sql.ErrNoRows:
			return genKey, nil
		case err != nil:
			return "", err
		}
	}
}

func SelectProjectCategories(project string) ([]Category, error) {

	sql := "select id,key,name,description from Category where project = ? and key != ?"
	rows, err := inst.Query(sql, project, DefaultCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cats := make([]Category, 0)
	for rows.Next() {
		cat := Category{}
		rows.Scan(&cat.Id, &cat.Key, &cat.Name, &cat.Description)
		cats = append(cats, cat)
	}
	return cats, nil
}

//TODO Tx
func DeleteCategory(catKey string) error {
	_, err := inst.Exec("delete from Category where key = ? ", catKey)
	return err
}

func SelectCategory(catId string) (Category, error) {
	cat := Category{}
	err := inst.QueryRow("select id,key,name,description from Category where key = ?", catId).
		Scan(&cat.Id, &cat.Key, &cat.Name, &cat.Description)
	return cat, err
}

func InsertCategory(cat Category) error {
	tx, err := inst.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = cat.Insert(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
