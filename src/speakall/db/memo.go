package db

import (
	"database/sql"
)

type Memo struct {
	Id      int
	Key     string
	Name    string
	Content string
}

func createMemoTable() error {
	_, err := Exec("CREATE TABLE Memo(id INTEGER PRIMARY KEY AUTOINCREMENT,key text,name text,content text)")
	return err
}

func deleteMemoTable() error {
	_, err := Exec("DROP TABLE if exists Memo")
	return err
}

func InsertMemo(key, name, content string) (sql.Result, error) {
	return inst.Exec("insert into Memo(key,name,content) values(?, ?, ?)", key, name, content)
}

func SelectMemo(key string) (Memo, error) {
	memo := Memo{}
	err := inst.QueryRow("select id,key,name,content from Memo where key = ?", key).
		Scan(&memo.Id, &memo.Key, &memo.Name, &memo.Content)
	return memo, err
}

func UpdateMemo(key, name, content string) error {
	_, err := inst.Exec("update memo set name=?,content=? where key = ?",
		name, content, key)
	return err
}

func DeleteMemo(key string) error {
	_, err := inst.Exec("delete from memo where key = ?",
		key)
	return err
}
