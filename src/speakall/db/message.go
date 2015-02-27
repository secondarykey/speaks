package db

import (
	"database/sql"
	"log"
)

type Message struct {
	Id       int
	Category string
	UserId   int
	Content  string
	Created  string
}

func createMessageTable() error {
	_, err := Exec("DROP TABLE if exists Message")
	if err != nil {
		return err
	}
	_, err = Exec("CREATE TABLE Message(id INTEGER PRIMARY KEY AUTOINCREMENT,category text,user_id integer,content text,created text)")
	return err
}

func SelectMessage(category, lastedId string) ([]Message, error) {

	sql := "select id,user_id,category,content,created from Message Where category = ?"
	if lastedId != "" {
		sql += " AND id < ?"
	} else {
		sql += " AND id < 9999999999"
	}
	sql += " ORDER BY created DESC LIMIT 10"

	log.Println(sql)

	rows, err := inst.Query(sql, category, lastedId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	msgs := make([]Message, 0)
	for rows.Next() {
		msg := Message{}
		rows.Scan(&msg.Id, &msg.UserId, &msg.Category,
			&msg.Content, &msg.Created)

		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func InsertMessage(userId int, category, content, created string) (sql.Result, error) {
	return inst.Exec("insert into Message(user_id,category,content,created) values(?, ?, ?, ?)", userId, category, content, created)
}
