package bm

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"
)

var (
	PartSeparator    = regexp.MustCompile("—+\\n")
	CommandSeparator = regexp.MustCompile("\\s+")
)

func TestFiles(t *testing.T) {
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

		before := parts[0]
		commands := CommandSeparator.Split(parts[1], -1)
		after := parts[2]

		temp, err := ioutil.TempFile("", file)
		if err != nil {
			t.Fatalf("cannot create temp file: %v", err)
		}
		defer os.Remove(temp.Name())
		t.Logf("using temp file %s", temp.Name())
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

		// interpret commands
		_ = commands

		data, err = ioutil.ReadFile(temp.Name())
		if err != nil {
			t.Fatalf("cannot read temp file: %v", err)
		}
		result := string(data)
		if result != after {
			t.Errorf("expected\n%s found\n%s", after, result)
		}
	}
}
