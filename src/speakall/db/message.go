package db

import (
	"database/sql"
	"log"
	"time"
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

func selectMessage(category string) ([]Message, error) {
	rows, err := inst.Query("select id,user_id,category,content,created from Message Where category = ?", category)
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

func InsertMessage(userId int, category string, content string) (sql.Result, error) {
	created := time.Now()
	countMessage()

	return inst.Exec("insert into Message(user_id,category,content,created) values(?, ?, ?, ?)", userId, category, content, created)
}

func countMessage() {
	row := inst.QueryRow("select count(*) from Message")
	var cnt int
	row.Scan(&cnt)
	log.Println(cnt)
}
