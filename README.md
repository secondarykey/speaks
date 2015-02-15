this document is still a draft.
becouse still in development,the documents is appropriate.

"SpeakAll" is Simple SNS.
use of on-premises environment.

environment :

  source activate
  (GOPATH + pwd)
  use thea 'deactive' command when ending.

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
- github.com/smartystreets/goconvey

run :
- go run ./src/*.go

test :
- goconvey(GOPATH/bin)


