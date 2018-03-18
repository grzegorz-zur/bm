package main

import (
	"flag"
	"fmt"
	bm "github.com/grzegorz-zur/bare-minimum"
	"log"
	"os"
)

func main() {
	logfile, err := os.Create("bm.log")
	if err != nil {
		fmt.Fprint(os.Stderr, "error opening logfile\n")
		os.Exit(-1)
	}
	log.SetOutput(logfile)
	defer logfile.Close()
	flag.Parse()
	base, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	paths := []string{}
	for _, path := range flag.Args() {
		paths = append(paths, path)
	}
	display := &bm.Display{}
	editor := bm.New(display, base, paths)
	editor.Start()
	editor.Wait()
}
