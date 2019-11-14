package main

import (
	"fmt"
	"github.com/rakyll/statik/fs"

	_ "github.com/secondarykey/speaks/statik"
)

func main() {

	fsys, err := fs.New()
	if err != nil {
		panic(err)
	}

	dir, err := fsys.Open("/")
	if err != nil {
		panic(err)
	}

	fis, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}

	for _, elm := range fis {
		fmt.Println(elm.Name())
	}

}
