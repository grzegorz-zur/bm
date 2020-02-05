package main

import (
	"fmt"
)

// Files represents all open files.
type Files struct {
	// File is the active file.
	*File
	list []*File
}

// Empty checks if no files are open.
func (fs *Files) Empty() bool {
	return len(fs.list) == 0
}

// Open opens file.
func (fs *Files) Open(path string) error {
	index, found := fs.find(path)
	if found {
		fs.switchFile(index)
		return nil
	}
	file, err := Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", file.Path, err)
	}
	index = fs.add(&file)
	fs.switchFile(index)
	return nil
}

// SwitchFile switches active file.
func (fs *Files) SwitchFile(d Direction) {
	if fs.Empty() {
		return
	}
	index := fs.current()
	index = wrap(index, len(fs.list), 1, d)
	fs.switchFile(index)
	fs.Read(false)
}

// Write writes all open files.
func (fs *Files) Write() error {
	if fs.Empty() {
		return nil
	}
	for _, file := range fs.list {
		_, err := file.Write()
		if err != nil {
			return fmt.Errorf("error writing file %s: %w", file.Path, err)
		}
	}
	return nil
}

// Close closes active file.
func (fs *Files) Close() {
	if fs.Empty() {
		return
	}
	index := fs.current()
	fs.remove(index)
	index = wrap(index, len(fs.list), 0, Forward)
	fs.switchFile(index)
}

func (fs *Files) add(f *File) int {
	fs.list = append(fs.list, f)
	return len(fs.list) - 1
}

func (fs *Files) remove(index int) {
	list := make([]*File, 0, len(fs.list)-1)
	list = append(list, fs.list[:index]...)
	list = append(list, fs.list[index+1:]...)
	fs.list = list
}

func (fs *Files) switchFile(index int) {
	if len(fs.list) != 0 {
		fs.File = fs.list[index]
	} else {
		fs.File = nil
	}
}

func (fs *Files) find(path string) (int, bool) {
	for i, file := range fs.list {
		if file.Path == path {
			return i, true
		}
	}
	return 0, false
}

func (fs *Files) current() int {
	for i, file := range fs.list {
		if fs.File == file {
			return i
		}
	}
	return 0
}
