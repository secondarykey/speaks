this document is still a draft.
becouse still in development,the documents is appropriate.

"speaks" is Simple SNS.
use of on-premises environment.

![TravisCI](https://travis-ci.org/secondarykey/speaks.svg?branch=master)

run :

go get -u github.com/secondarykey/speaks
go install

speaks init
speaks start

environment :

  .speaks 

default setting :

- http://localhost:5555
- admin@localhost/p@ssword

```

go-bindata.exe -pkg=config -o=./config/binary.go .speaks/...
```


not yet Modified Template for 0.5.0

database.tmpl
category.tmpl
user.tmpl
memo/
  edit.tmpl
  list.tmpl
  view.tmpl

