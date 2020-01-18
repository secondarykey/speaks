package db

import (
	"database/sql"
	"fmt"
)

type Message struct {
	Id       int
	Category string
	Project  string
	UserId   int
	UserName string
	Content  string
	Created  string
}

func (m Message) Create(tx *sql.Tx) error {
	_, err := Exec("CREATE TABLE Message(id INTEGER PRIMARY KEY AUTOINCREMENT,project text,category text,user_id integer,content text,created text)")
	return err
}

func (m Message) Drop(tx *sql.Tx) error {
	_, err := Exec("DROP TABLE if exists Message")
	return err
}

func (m Message) Init(tx *sql.Tx) error {
	return fmt.Errorf("Not yet Implemented")
}

func (m Message) Insert(tx *sql.Tx) (sql.Result, error) {

	sql := "INSERT INTO Message(user_id,project,category,content,created) values(?, ?, ?, ?, ?)"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(m.Id, m.Project, m.Category, m.Content, m.Created)
}

func createMessageTable() error {
	msg := Message{}
	return msg.Create(nil)
}

func dropMessageTable() error {
	msg := Message{}
	return msg.Drop(nil)
}

func SelectMessages(project, category, lastedId string) ([]Message, error) {

	sql := createMessageSQL()
	if lastedId != "" {
		sql += " AND Message.id < ?"
	} else {
		sql += " AND Message.id < 9999999999"
	}
	sql += " ORDER BY Message.created DESC LIMIT 10"

	rows, err := inst.Query(sql, category, project, lastedId)
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

func SelectAllMessage(category, project string) ([]Message, error) {
	sql := createMessageSQL()
	sql += " ORDER BY Message.created ASC"
	rows, err := inst.Query(sql, category, project)
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

func SearchMessages(project, category, text string, page int) ([]Message, error) {

	limit := 10
	offset := (page - 1) * limit

	sql := createSearchMessageSQL()
	if offset <= 0 {
		sql += fmt.Sprintf(" ORDER BY Message.created DESC LIMIT %d", limit)
	} else {
		sql += fmt.Sprintf(" ORDER BY Message.created DESC LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := inst.Query(sql, category, project, "%"+text+"%")
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

func createMessageSQL() string {
	sql := "select " +
		"Message.id," +
		"Message.user_id," +
		"Message.category," +
		"Message.content," +
		"Message.created," +
		"User.Name" +
		" from Message INNER JOIN User ON Message.user_id = User.id" +
		" Where category = ? and project = ?"
	return sql
}

func createSearchMessageSQL() string {
	sql := "select " +
		"Message.id," +
		"Message.user_id," +
		"Message.category," +
		"Message.content," +
		"Message.created," +
		"User.Name" +
		" from Message INNER JOIN User ON Message.user_id = User.id" +
		" Where category = ? and project = ? and content like ?"
	return sql
}

func InsertMessage(userId int, project, category, content, created string) error {

	tx, err := inst.Begin()
	if err != nil {
		return err
	}

	msg := Message{}
	msg.Id = userId
	msg.Project = project
	msg.Category = category
	msg.Content = content
	msg.Created = created

	_, err = msg.Insert(tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func DeleteMessage(id int) error {

	tx, err := inst.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sql := "delete from Message where id = ?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteAllMessage(tx *sql.Tx, category string) (sql.Result, error) {
	stmt, err := tx.Prepare("delete from Message where category = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	return stmt.Exec(category)
}

func createMemoContent(key, project string) (string, error) {

	msgs, err := SelectAllMessage(key, project)
	if err != nil {
		return "", err
	}

	content := ""
	for _, elm := range msgs {
		content += elm.UserName + ":" + elm.Created
		content += "\n"
		content += elm.Content
		content += "\n\n"
	}
	return content, nil
}
