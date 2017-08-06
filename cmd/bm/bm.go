package main

import (
	bm "github.com/grzegorz-zur/bare-minimum"
	"log"
)

func main() {

	editor, err := bm.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = editor.Run()
	if err != nil {
		log.Fatal(err)
	}

}
