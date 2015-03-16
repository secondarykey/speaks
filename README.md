this document is still a draft.
becouse still in development,the documents is appropriate.

"SpeakAll" is Simple SNS.
use of on-premises environment.

![TravisCI](https://travis-ci.org/secondarykey/SpeakAll.svg?branch=master)

environment :

  put the SpeakAll.ini the run dhirectory.

using (go get):
- golang.org/x/net/websocket
- github.com/gorilla/sessions
- github.com/mattn/go-sqlite3
- github.com/satori/go.uuid
- github.com/BurntSushi/toml
- github.com/smartystreets/goconvey

default setting :

- http://localhost:5555
- admin@localhost/password

develop :

  source activate
  add GOPATH(root dhirectory)
  use the 'deactive' command when ending.

- go run ./src/main.go

test :
- goconvey(GOPATH/bin)
