// Bare Minimum.
//
// Minimalistic text editor.
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
	user := cu.Username
	date := time.Now().Format("2006-01-02_15-04-05")
	prefix := fmt.Sprintf("bm_%s_%s_", user, date)

	lf, err := ioutil.TempFile("", prefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening log file %v", err)
		os.Exit(-1)
	}
	defer func() {
		err = lf.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error closing log file %v", err)
		}
	}()
	log.SetOutput(lf)
	defer log.SetOutput(os.Stderr)

	flag.Parse()
	ps := []string{}
	for _, p := range flag.Args() {
		ps = append(ps, p)
	}

	d := &Display{}
	e := New(d, ps)
	e.Start()
	e.Wait()

}
