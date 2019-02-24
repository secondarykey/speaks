this document is still a draft.
becouse still in development,the documents is appropriate.

"speaks" is Simple SNS.
use of on-premises environment.

![TravisCI](https://travis-ci.org/secondarykey/SpeakAll.svg?branch=master)

environment :

  put the speaks.ini the run directory.

default setting :

- http://localhost:5555
- admin@localhost/password

run :

```
package main

import (
    "fmt"
    "os"

    "github.com/secondarykey/speaks"
)

func main() {
    err = speaks.Listen()
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    os.Exit(0)
}
```

- go run main.go

test :

none test
