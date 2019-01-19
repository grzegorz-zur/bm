package bm

import (
	"bufio"
	"bytes"
	"errors"
	tb "github.com/nsf/termbox-go"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
	"unicode"
)

func TestEditor(t *testing.T) {
	prefix := time.Now().Format("bm_2006-01-02_15-04-05_")
	temp, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("error creating temporary directory %s: %v", prefix, err)
	}
	names, err := list("test")
	if err != nil {
		t.Fatalf("error reading test directory: %v", err)
	}
	for _, name := range names {
		test := test(name, temp)
		t.Run(name, test)
	}
}

func test(name, temp string) func(t *testing.T) {
	return func(t *testing.T) {
		work, files, err := setup(name, temp, t)
		if err != nil {
			t.Fatalf("error on setup: %s: %v", name, err)
		}
		cmds, err := commands(name)
		if err != nil {
			t.Fatalf("error parsing script: %s: %v", name, err)
		}

		dir, err := os.Getwd()
		if err != nil {
			t.Fatalf("error getting current directory: %v", err)
		}
		defer os.Chdir(dir)
		err = os.Chdir(work)
		if err != nil {
			t.Fatalf("error changing dir to %s: %v", work, err)
		}

		editor := New(nil, files)
		editor.Start()

		err = interpret(editor, cmds)
		if err != nil {
			t.Fatalf("error interpreting script: %s: %v", name, err)
		}

		editor.Wait()

		os.Chdir(dir)
		if err != nil {
			t.Fatalf("error changing dir to %s: %v", dir, err)
		}
		err = verify(name, work, t)
		if err != nil {
			t.Fatalf("wrong results: %s: %v", name, err)
		}
	}
}

func setup(name, temp string, t *testing.T) (work string, files []string, err error) {
	in := path.Join("test", name, "in")
	files, err = list(in)
	if err != nil && !os.IsNotExist(err) {
		return "", nil, err
	}
	work = path.Join(temp, name)
	err = os.MkdirAll(work, os.ModePerm)
	if err != nil {
		return "", nil, err
	}
	for _, file := range files {
		src := path.Join(in, file)
		dst := path.Join(work, file)
		err := copy(src, dst)
		if err != nil {
			return "", nil, err
		}
	}
	return
}

func commands(name string) (cmds []string, err error) {
	path := path.Join("test", name, "script")
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return
		}
		cmd := scanner.Text()
		cmds = append(cmds, cmd)
	}
	return
}

func verify(name, work string, t *testing.T) (err error) {
	out := path.Join("test", name, "out")
	expected, err := list(out)
	if err != nil {
		return
	}
	actual, err := list(work)
	if err != nil {
		return
	}
	if len(expected) != len(actual) {
		t.Logf("expected files %v", expected)
		t.Logf("actual files %v", actual)
		t.FailNow()
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Logf("expected files %v", expected)
			t.Logf("actual files %v", actual)
			t.FailNow()
		}
	}
	for _, exp := range expected {
		actualPath := path.Join(work, exp)
		actualContent, err := ioutil.ReadFile(actualPath)
		if err != nil {
			return err
		}
		expectedPath := path.Join("test", name, "out", exp)
		expectedContent, err := ioutil.ReadFile(expectedPath)
		if err != nil {
			return err
		}
		if bytes.Compare(actualContent, expectedContent) != 0 {
			t.Log("comparing file " + exp)
			t.Log("expected content")
			t.Log(string(expectedContent))
			t.Log("actual content")
			t.Log(string(actualContent))
			t.FailNow()
		}
	}
	return
}

func interpret(editor *Editor, commands []string) (err error) {
	for _, cmd := range commands {
		runes := []rune(cmd)
		switch {
		case len(cmd) == 1:
			editor.keys <- tb.Event{Ch: runes[0]}
		case len(cmd) == 2 && runes[0] == '^':
			letter := unicode.ToUpper(runes[1])
			offset := int(letter - 'A')
			key := tb.KeyCtrlA + tb.Key(offset)
			editor.keys <- tb.Event{Key: key}
		case cmd == "escape":
			editor.keys <- tb.Event{Key: tb.KeyEsc}
		case cmd == "left":
			editor.keys <- tb.Event{Key: tb.KeyArrowLeft}
		case cmd == "right":
			editor.keys <- tb.Event{Key: tb.KeyArrowRight}
		case cmd == "up":
			editor.keys <- tb.Event{Key: tb.KeyArrowUp}
		case cmd == "down":
			editor.keys <- tb.Event{Key: tb.KeyArrowDown}
		case cmd == "space":
			editor.keys <- tb.Event{Key: tb.KeySpace}
		case cmd == "tab":
			editor.keys <- tb.Event{Key: tb.KeyTab}
		case cmd == "enter":
			editor.keys <- tb.Event{Key: tb.KeyEnter}
		case cmd == "backspace":
			editor.keys <- tb.Event{Key: tb.KeyBackspace2}
		case cmd == "delete":
			editor.keys <- tb.Event{Key: tb.KeyDelete}
		case cmd == "TOUCH":
			t := time.Now().Local()
			err = os.Chtimes(editor.File.Path, t, t)
			if err != nil {
				return
			}
		case cmd == "CHECK":
			editor.Check()
		default:
			err = errors.New("unknown command: " + cmd)
			return
		}
	}
	return
}

func list(path string) (names []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	for _, file := range files {
		name := file.Name()
		names = append(names, name)
	}
	return
}

func copy(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
