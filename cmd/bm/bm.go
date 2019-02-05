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
		fmt.Fprintf(os.Stderr, "error opening log file %v", err)
		os.Exit(-1)
	}
	defer func() {
		err = logfile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error closing log file %v", err)
		}
	}()
	log.SetOutput(logfile)
	defer log.SetOutput(os.Stderr)
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
