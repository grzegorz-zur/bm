package main

import (
	"flag"
	"fmt"
	bm "github.com/grzegorz-zur/bare-minimum"
	"log"
	"os"
)

func main() {
	l, err := os.Create("bm.log")
	if err != nil {
		fmt.Fprint(os.Stderr, "cannot initiate log\n")
		os.Exit(-1)
	}
	log.SetOutput(l)
	defer l.Close()
	flag.Parse()
	editor := bm.New()
	for _, path := range flag.Args() {
		err = editor.Open(path)
		if err != nil {
			log.Println(err)
		}
	}
	err = editor.Run()
	if err != nil {
		log.Fatal(err)
	}
}
