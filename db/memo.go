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
	_, err := Exec("CREATE TABLE Memo(id INTEGER PRIMARY KEY AUTOINCREMENT,key text,project text,name text,content text)")
	return err
}

func dropMemoTable() error {
	_, err := Exec("DROP TABLE if exists Memo")
	return err
}

func InsertMemo(tx *sql.Tx, key, project, name, content string) (sql.Result, error) {
	//tx
	stmt, err := tx.Prepare("insert into Memo(key,project,name,content) values(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(key, project, name, content)

}

func SelectMemo(key, project string) (Memo, error) {
	memo := Memo{}
	err := inst.QueryRow("select id,key,name,content from Memo where key = ? and project = ?", key, project).
		Scan(&memo.Id, &memo.Key, &memo.Name, &memo.Content)
	return memo, err
}

func SelectProjectMemo(project string) ([]Memo, error) {

	sql := "SELECT id,key,name,content from Memo"
	sql += " WHERE project = ?"

	rows, err := inst.Query(sql, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	memos := make([]Memo, 0)
	for rows.Next() {
		memo := Memo{}
		rows.Scan(&memo.Id, &memo.Key, &memo.Name, &memo.Content)
		memos = append(memos, memo)
	}
	return memos, nil
}

func UpdateMemo(key, project, name, content string) error {
	_, err := inst.Exec("update memo set name=?,content=? where key = ? and project = ?",
		name, content, key, project)
	return err
}

func DeleteMemo(key, project string) error {
	_, err := inst.Exec("delete from memo where key = ? and project = ?",
		key, project)
	return err
}
