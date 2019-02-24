package db

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
)

type Category struct {
	Id          int
	Key         string
	Name        string
	Description string
}

func createCategoryTable() error {
	_, err := Exec("CREATE TABLE Category(id INTEGER PRIMARY KEY AUTOINCREMENT,key text,name text,description text)")
	return err
}

func deleteCategoryTable() error {
	_, err := Exec("DROP TABLE if exists Category")
	return err
}

func InsertCategory(key, name, desc string) (sql.Result, error) {
	return inst.Exec("insert into Category(key,name,description) values(?, ?, ?)", key, name, desc)
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

func SelectAllCategory() ([]Category, error) {
	sql := "select id,key,name,description from Category"
	rows, err := inst.Query(sql)
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

func DeleteCategory(catId string) error {
	_, err := inst.Exec("delete from Category where key = ? ", catId)
	return err
}

func SelectCategory(catId string) (Category, error) {
	cat := Category{}
	err := inst.QueryRow("select id,key,name,description from Category where key = ?", catId).
		Scan(&cat.Id, &cat.Key, &cat.Name, &cat.Description)
	return cat, err
}
