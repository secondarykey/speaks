environment :

  source activate

initialize :

- http://localhost:5555
- admin@localhost/password

loading SpeakAll.ini

using (go get):
- golang.org/x/net/websocket
- github.com/gorilla/sessions
- github.com/mattn/go-sqlite3
- github.com/satori/go.uuid
- github.com/BurntSushi/toml

run :
- go run ./src/main.go
