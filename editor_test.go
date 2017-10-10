package bm

import (
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
)

var (
	PartSeparator    = regexp.MustCompile("\\n?—+\\n")
	CommandSeparator = regexp.MustCompile("\\s+")
)

func TestEditor(t *testing.T) {
	files := list(t, "test")
	for _, file := range files {
		t.Run(file, testCase(file))
	}
}

func list(t *testing.T, dir string) (names []string) {
	t.Helper()
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("cannot read directory", err)
	}
	for _, fi := range fis {
		if fi.Mode().IsRegular() {
			names = append(names, fi.Name())
		}
	}
	return
}

func testCase(file string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Logf("testing file %s", file)

		path := path.Join("test", file)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("cannot read test file: %v", err)
		}

		text := string(data)
		parts := PartSeparator.Split(text, -1)
		if len(parts) != 3 {
			t.Fatalf("invalid test file %s", file)
		}
		t.Logf("parts %v", parts)

		before := parts[0]
		commands := []string{}
		if strings.Trim(parts[1], " ") != "" {
			commands = CommandSeparator.Split(parts[1], -1)
		}
		after := parts[2]

		temp, err := ioutil.TempFile("", file+"_")
		if err != nil {
			t.Fatalf("cannot create temp file: %v", err)
		}
		defer os.Remove(temp.Name())
		t.Logf("temp file %s", temp.Name())
		_, err = temp.WriteString(before)
		if err != nil {
			t.Fatalf("cannot write to temp file: %v", err)
		}
		err = temp.Close()
		if err != nil {
			t.Fatalf("cannot close temp file: %v", err)
		}
		editor, err := Open(temp.Name())
		if err != nil {
			t.Fatalf("cannot open editor: %v", err)
		}
		defer editor.Close()

		t.Logf("running commands: %v", commands)
		err = interpret(t, editor, commands)
		if err != nil {
			t.Fatal(err)
		}

		err = editor.Write()
		if err != nil {
			t.Fatal(err)
		}

		data, err = ioutil.ReadFile(temp.Name())
		if err != nil {
			t.Fatalf("cannot read temp file: %v", err)
		}
		result := string(data)
		if result != after {
			t.Errorf("comparison failed\nexpected\n\"%s\"\nfound\n\"%s\"", after, result)
		}
	}
}

func interpret(t *testing.T, editor *Editor, commands []string) (err error) {
	for _, cmd := range commands {
		t.Logf("running command: \"%s\"", cmd)
		switch {
		case len(cmd) == 1:
			runes := []rune(cmd)
			event := tb.Event{Ch: runes[0]}
			err = editor.Key(event)
		default:
			err = errors.Errorf("cannot interpret command \"%s\"", cmd)
		}
	}
	return
}
