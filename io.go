package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Open opens a file.
func Open(path string) (file File, err error) {
	path = filepath.Clean(path)
	file = File{
		Path:    path,
		history: &History{},
	}
	_, err = file.Read(false)
	if err != nil {
		return file, fmt.Errorf("error opening file %s: %w", path, err)
	}
	return file, nil
}

// Read loads thie file.
func (file *File) Read(force bool) (read bool, err error) {
	stat, err := os.Stat(file.Path)
	if os.IsNotExist(err) {
		err = ioutil.WriteFile(file.Path, []byte{}, os.ModePerm)
		if err != nil {
			return false, fmt.Errorf("error creating file %s info: %w", file.Path, err)
		}
		stat, err = os.Stat(file.Path)
	}
	if err != nil {
		return false, fmt.Errorf("error checking file %s info: %w", file.Path, err)
	}
	changed := file.time != stat.ModTime()
	if !changed && !force {
		return false, nil
	}
	data, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return false, fmt.Errorf("error reading file %s: %w", file.Path, err)
	}
	file.content = string(data)
	stat, err = os.Stat(file.Path)
	if err != nil {
		return true, fmt.Errorf("error checking file %s info: %w", file.Path, err)
	}
	file.changed = false
	file.time = stat.ModTime()
	file.Archive()
	if file.location > len(file.content) {
		file.location = len(file.content)
	}
	if file.mark > len(file.content) {
		file.mark = len(file.content)
	}
	return true, nil
}

// Write writes file contents.
func (file *File) Write() (wrote bool, err error) {
	stat, err := os.Stat(file.Path)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("error checking file %s info: %w", file.Path, err)
	}
	exists := stat != nil
	changed := true
	if exists {
		changed = file.time != stat.ModTime()
	}
	if !file.changed && !changed {
		return false, nil
	}
	data := []byte(file.content)
	err = ioutil.WriteFile(file.Path, data, os.ModePerm)
	if err != nil {
		return false, fmt.Errorf("error writing file %s: %w", file.Path, err)
	}
	stat, err = os.Stat(file.Path)
	if err != nil {
		return true, fmt.Errorf("error checking file %s info: %w", file.Path, err)
	}
	file.changed = false
	file.time = stat.ModTime()
	return true, nil
}
