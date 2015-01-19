package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var inst *sql.DB

func Listen(path string) {
	var err error
	inst, err = sql.Open("sqlite3", path)
	if err != nil {
		panic(err.Error())
	}
}

/* create
sqls := []string{
	"create table foo (id integer not null primary key, name text)",
	"delete from foo",
}
for _, sql := range sqls {
	_, err = db.Exec(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		return
	}
}
*/

/* tx insert
tx, err := db.Begin()
if err != nil {
	fmt.Println(err)
	return
}
stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
if err != nil {
	fmt.Println(err)
	return
}
defer stmt.Close()

for i := 0; i < 100; i++ {
	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
	if err != nil {
		fmt.Println(err)
		return
	}
}
tx.Commit()
*/

func Select() {

	rows, err := inst.Query("select id, name from foo")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		println(id, name)
	}
}
