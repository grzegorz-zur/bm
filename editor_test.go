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
	Terminator = regexp.MustCompile("\n?—+\n")
	WhiteChars = regexp.MustCompile("\\s+")
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
		parts := Terminator.Split(text, -1)
		if len(parts) != 4 {
			t.Fatalf("invalid test file %s", file)
		}

		before := parts[0]
		commands := []string{}
		if strings.Trim(parts[1], " ") != "" {
			commands = WhiteChars.Split(parts[1], -1)
		}
		after := parts[2]
		t.Logf("before\n`%s`", before)
		t.Logf("commands: %v", commands)
		t.Logf("after\n`%s`", after)

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
		editor := New()
		err = editor.Open(temp.Name())
		if err != nil {
			t.Fatalf("cannot open editor: %v", err)
		}
		defer editor.Close()

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
			t.Errorf("comparison failed")
			t.Errorf("expected\n%s", after)
			t.Errorf("found\n%s", result)
		}
	}
}

func interpret(t *testing.T, editor *Editor, commands []string) (err error) {
	for _, cmd := range commands {
		t.Logf("event: %s", cmd)
		var event tb.Event
		switch {
		case len(cmd) == 1:
			runes := []rune(cmd)
			event = tb.Event{Ch: runes[0]}
		case cmd == "escape":
			event = tb.Event{Key: tb.KeyEsc}
		case cmd == "left":
			event = tb.Event{Key: tb.KeyArrowLeft}
		case cmd == "right":
			event = tb.Event{Key: tb.KeyArrowRight}
		case cmd == "up":
			event = tb.Event{Key: tb.KeyArrowUp}
		case cmd == "down":
			event = tb.Event{Key: tb.KeyArrowDown}
		case cmd == "space":
			event = tb.Event{Key: tb.KeySpace}
		case cmd == "tab":
			event = tb.Event{Key: tb.KeyTab}
		case cmd == "enter":
			event = tb.Event{Key: tb.KeyEnter}
		case cmd == "backspace":
			event = tb.Event{Key: tb.KeyBackspace}
		case cmd == "delete":
			event = tb.Event{Key: tb.KeyDelete}
		default:
			err = errors.Errorf("cannot interpret command %s", cmd)
			return
		}
		err = editor.Key(event)
	}
	return
}
