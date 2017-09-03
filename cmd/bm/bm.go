package main

import (
	"flag"
	"fmt"
	bm "github.com/grzegorz-zur/bare-minimum"
	"log"
	"os"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		fmt.Fprint(os.Stderr, "no file name given\n")
		flag.Usage()
		os.Exit(-1)
	}
	editor, err := bm.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	err = editor.Run()
	if err != nil {
		log.Fatal(err)
	}
}
