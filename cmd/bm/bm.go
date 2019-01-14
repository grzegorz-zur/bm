package main

import (
	"flag"
	"fmt"
	bm "github.com/grzegorz-zur/bare-minimum"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	logfile, err := ioutil.TempFile("", "bm")
	if err != nil {
		fmt.Fprint(os.Stderr, "error opening logfile\n")
		os.Exit(-1)
	}
	log.SetOutput(logfile)
	defer logfile.Close()
	flag.Parse()
	paths := []string{}
	for _, path := range flag.Args() {
		paths = append(paths, path)
	}
	display := &bm.Display{}
	editor := bm.New(display, paths)
	editor.Start()
	editor.Wait()
}
