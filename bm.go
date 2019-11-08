package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"
)

func main() {
	cu, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current user %v", err)
		os.Exit(-1)
	}
	username := cu.Username
	date := time.Now().Format("2006-01-02_15-04-05")
	prefix := fmt.Sprintf("bm_%s_%s_", username, date)
	logfile, err := ioutil.TempFile("", prefix)
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
	display := &Display{}
	editor := New(display, paths)
	editor.Start()
	editor.Wait()
}
