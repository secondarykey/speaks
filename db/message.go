package db

import (
	"database/sql"
)

type Message struct {
	Id       int
	Category string
	UserId   int
	UserName string
	Content  string
	Created  string
}

func createMessageTable() error {
	_, err := Exec("CREATE TABLE Message(id INTEGER PRIMARY KEY AUTOINCREMENT,category text,user_id integer,content text,created text)")
	return err
}

func deleteMessageTable() error {
	_, err := Exec("DROP TABLE if exists Message")
	return err
}

func SelectMessage(category, lastedId string) ([]Message, error) {

	sql := createSQL()
	if lastedId != "" {
		sql += " AND Message.id < ?"
	} else {
		sql += " AND Message.id < 9999999999"
	}
	sql += " ORDER BY Message.created DESC LIMIT 10"

	rows, err := inst.Query(sql, category, lastedId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	msgs := make([]Message, 0)
	for rows.Next() {
		msg := Message{}
		rows.Scan(&msg.Id, &msg.UserId, &msg.Category,
			&msg.Content, &msg.Created, &msg.UserName)

		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func SelectAllMessage(category string) ([]Message, error) {
	sql := createSQL()
	sql += " ORDER BY Message.created ASC"
	rows, err := inst.Query(sql, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	msgs := make([]Message, 0)
	for rows.Next() {
		msg := Message{}
		rows.Scan(&msg.Id, &msg.UserId, &msg.Category,
			&msg.Content, &msg.Created, &msg.UserName)

		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func createSQL() string {
	sql := "select " +
		"Message.id," +
		"Message.user_id," +
		"Message.category," +
		"Message.content," +
		"Message.created," +
		"User.Name" +
		" from Message INNER JOIN User ON Message.user_id = User.id" +
		" Where category = ?"
	return sql
}

func InsertMessage(userId int, category, content, created string) (sql.Result, error) {
	return inst.Exec("insert into Message(user_id,category,content,created) values(?, ?, ?, ?)", userId, category, content, created)
}

func DeleteMessage(id string, user_id int) error {
	_, err := inst.Exec("delete from Message where id = ? and user_id = ?",
		id, user_id)
	return err
}

func DeleteAllMessage(category string) error {
	_, err := inst.Exec("delete from Message where category = ? ", category)
	return err
}
