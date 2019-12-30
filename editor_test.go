package main

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/gdamore/tcell"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
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
			t.Fatalf("error parsing script %s: %v", name, err)
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

		logpath := path.Join(temp, name+".log")
		logfile, err := os.Create(logpath)
		if err != nil {
			t.Fatalf("error opening logfile %s: %v", logpath, err)
			os.Exit(-1)
		}
		defer logfile.Close()
		log.SetOutput(logfile)
		defer log.SetOutput(os.Stderr)

		newScreen := func() (tcell.Screen, error) {
			return tcell.NewSimulationScreen(""), nil
		}
		editor := New(newScreen, files)
		err = editor.Start()
		if err != nil {
			t.Fatalf("error starting editor: %v", err)
		}

		err = interpret(editor, cmds)
		if err != nil {
			t.Fatalf("error interpreting script: %s: %v", name, err)
		}

		editor.Wait()

		logfile.Close()
		log.SetOutput(os.Stderr)

		err = os.Chdir(dir)
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
			t.Log(format(expectedContent))
			t.Log("actual content")
			t.Log(format(actualContent))
			t.FailNow()
		}
	}
	return
}

func format(content []byte) (text string) {
	text = string(content)
	text = strings.Replace(text, " ", "␣", -1)
	text = strings.Replace(text, "\t", "␉", -1)
	text = strings.Replace(text, "\r", "␍", -1)
	text = strings.Replace(text, "\n", "␊", -1)
	return
}

func interpret(editor *Editor, commands []string) (err error) {
	for _, cmd := range commands {
		runes := []rune(cmd)
		switch {
		case len(cmd) == 1:
			rune := runes[0]
			sendRune(editor, rune)
		case len(cmd) == 2 && runes[0] == '^':
			letter := unicode.ToUpper(runes[1])
			offset := int(letter - 'A')
			key := tcell.KeyCtrlA + tcell.Key(offset)
			sendKey(editor, key)
		case cmd == "left":
			sendKey(editor, tcell.KeyLeft)
		case cmd == "right":
			sendKey(editor, tcell.KeyRight)
		case cmd == "up":
			sendKey(editor, tcell.KeyUp)
		case cmd == "down":
			sendKey(editor, tcell.KeyDown)
		case cmd == "space":
			sendRune(editor, ' ')
		case cmd == "tab":
			sendKey(editor, tcell.KeyTab)
		case cmd == "enter":
			sendKey(editor, tcell.KeyEnter)
		case cmd == "backspace":
			sendKey(editor, tcell.KeyBackspace2)
		case cmd == "delete":
			sendKey(editor, tcell.KeyDelete)
		case cmd == "ctrlspace":
			sendKey(editor, tcell.KeyCtrlSpace)
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

func sendKey(editor *Editor, key tcell.Key) {
	ev := tcell.NewEventKey(key, 0, 0)
	editor.events <- ev
}

func sendRune(editor *Editor, r rune) {
	ev := tcell.NewEventKey(tcell.KeyRune, r, 0)
	editor.events <- ev
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
