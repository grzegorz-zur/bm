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

const (
	base   = "test"
	in     = "in"
	out    = "out"
	script = "script"
)

func TestEditor(t *testing.T) {
	prefix := time.Now().Format("bm_2006-01-02_15-04-05_")
	temp, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("cannot create temporary directory %s: %v", prefix, err)
	}
	files, err := ioutil.ReadDir(base)
	if err != nil {
		t.Fatalf("cannot read test directory %s: %v", base, err)
	}
	for _, file := range files {
		name := file.Name()
		test := test(name, base, temp)
		t.Run(name, test)
	}
}

func test(name, base, temp string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		err := setup(name, base, temp)
		if err != nil {
			t.Fatalf("setup failure: %v", err)
		}
		cmds, err := commands(name, base)
		if err != nil {
			t.Fatalf("script failure: %v", err)
		}

		testPath := path.Join(temp, name)
		editor := New(testPath)
		defer editor.Close()

		files, err := ioutil.ReadDir(testPath)
		for _, file := range files {
			err := editor.Open(file.Name())
			if err != nil {
				t.Fatalf("cannot open editor: %v", err)
			}
		}
		editor.Next(Forward)

		err = interpret(editor, cmds)
		if err != nil {
			t.Fatalf("cannot write files: %v", err)
		}

		err = editor.WriteAll()
		if err != nil {
			t.Fatalf("cannot write files: %v", err)
		}

		err = verify(name, base, temp, t)
		if err != nil {
			t.Fatalf("cannot verify files: %v", err)
		}
	}
}

func setup(name, base, temp string) (err error) {
	inPath := path.Join(base, name, in)
	files, err := ioutil.ReadDir(inPath)
	if err != nil {
		return err
	}
	dir := path.Join(temp, name)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	for _, file := range files {
		src := path.Join(inPath, file.Name())
		dst := path.Join(dir, file.Name())
		err := copy(src, dst)
		if err != nil {
			return err
		}
	}
	return
}

func commands(name, base string) (cmds []string, err error) {
	scriptPath := path.Join(base, name, script)
	file, err := os.Open(scriptPath)
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

func verify(name, base, temp string, t *testing.T) (err error) {
	outPath := path.Join(base, name, out)
	expected, err := list(outPath)
	if err != nil {
		return
	}
	testPath := path.Join(temp, name)
	actual, err := list(testPath)
	if err != nil {
		return
	}
	if len(expected) != len(actual) {
		t.Logf("expected files %v", expected)
		t.Logf("actual files %v", actual)
		t.Fail()
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Logf("expected files %v", expected)
			t.Logf("actual files %v", actual)
			t.Fail()
		}
	}
	for i := range expected {
		actualPath := path.Join(temp, name, expected[i])
		actualContent, err := ioutil.ReadFile(actualPath)
		if err != nil {
			return err
		}
		expectedPath := path.Join(base, name, out, expected[i])
		expectedContent, err := ioutil.ReadFile(expectedPath)
		if err != nil {
			return err
		}
		if bytes.Compare(actualContent, expectedContent) != 0 {
			t.Log("expected content")
			t.Log(string(expectedContent))
			t.Log("actual content")
			t.Log(string(actualContent))
			t.Fail()
		}
	}
	return
}

func interpret(editor *Editor, commands []string) (err error) {
	for _, cmd := range commands {
		runes := []rune(cmd)
		var event tb.Event
		switch {
		case len(cmd) == 1:
			event = tb.Event{Ch: runes[0]}
		case len(cmd) == 2 && runes[0] == '^':
			letter := unicode.ToUpper(runes[1])
			offset := int(letter - 'A')
			key := tb.KeyCtrlA + tb.Key(offset)
			event = tb.Event{Key: key}
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
			event = tb.Event{Key: tb.KeyBackspace2}
		case cmd == "delete":
			event = tb.Event{Key: tb.KeyDelete}
		default:
			err = errors.New("cannot interpret command: " + cmd)
			return
		}
		err = editor.Key(event)
		if err != nil {
			return err
		}
	}
	return
}

func list(path string) (names []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
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
