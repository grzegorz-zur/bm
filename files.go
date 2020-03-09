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
func (files *Files) Empty() bool {
	return len(files.list) == 0
}

// Open opens file.
func (files *Files) Open(path string) error {
	position, found := files.find(path)
	if found {
		files.switchFile(position)
		return nil
	}
	file, err := Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", file.Path, err)
	}
	position = files.add(&file)
	files.switchFile(position)
	return nil
}

// SwitchFile switches active file.
func (files *Files) SwitchFile(direction Direction) {
	if files.Empty() {
		return
	}
	position := files.current()
	position = wrap(position, len(files.list), 1, direction)
	files.switchFile(position)
	files.Read(false)
}

// Write writes all open files.
func (files *Files) Write() error {
	if files.Empty() {
		return nil
	}
	for _, file := range files.list {
		_, err := file.Write()
		if err != nil {
			return fmt.Errorf("error writing file %s: %w", file.Path, err)
		}
	}
	return nil
}

// Close closes active file.
func (files *Files) Close() {
	if files.Empty() {
		return
	}
	position := files.current()
	files.remove(position)
	if files.Empty() {
		position = 0
		return
	}
	position = wrap(position, len(files.list), 0, Forward)
	files.switchFile(position)
}

func (files *Files) add(file *File) int {
	files.list = append(files.list, file)
	return len(files.list) - 1
}

func (files *Files) remove(position int) {
	list := make([]*File, 0, len(files.list)-1)
	list = append(list, files.list[:position]...)
	list = append(list, files.list[position+1:]...)
	files.list = list
}

func (files *Files) switchFile(position int) {
	if len(files.list) != 0 {
		files.File = files.list[position]
	} else {
		files.File = nil
	}
}

func (files *Files) find(path string) (position int, ok bool) {
	for position, file := range files.list {
		if file.Path == path {
			return position, true
		}
	}
	return 0, false
}

func (files *Files) current() int {
	for position, file := range files.list {
		if files.File == file {
			return position
		}
	}
	return 0
}
