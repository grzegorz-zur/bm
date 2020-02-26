package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	perm os.FileMode = 0644
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
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("error reading file %s info: %w", file.Path, err)
	}
	exists := stat != nil
	changed := true
	if exists {
		changed = file.time != stat.ModTime()
	}
	if !changed && !force {
		return false, nil
	}
	flags := os.O_RDWR | os.O_CREATE
	osFile, err := os.OpenFile(file.Path, flags, perm)
	if err != nil {
		return false, fmt.Errorf("error opening file %s: %w", file.Path, err)
	}
	defer osFile.Close()
	buffer := &bytes.Buffer{}
	io.Copy(buffer, osFile)
	file.content = string(buffer.Bytes())
	osFile.Close()
	stat, err = os.Stat(file.Path)
	if err != nil {
		return true, fmt.Errorf("error reading file %s info: %w", file.Path, err)
	}
	file.changed = false
	file.time = stat.ModTime()
	file.Archive()
	return true, nil
}

// Write writes file contents.
func (file *File) Write() (wrote bool, err error) {
	stat, err := os.Stat(file.Path)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("error reading file %s info: %w", file.Path, err)
	}
	exists := stat != nil
	changed := true
	if exists {
		changed = file.time != stat.ModTime()
	}
	if !file.changed && !changed {
		return false, nil
	}
	osFile, err := os.Create(file.Path)
	defer osFile.Close()
	if err != nil {
		return false, fmt.Errorf("error writing file %s: %w", file.Path, err)
	}
	bytes := []byte(file.content)
	osFile.Write(bytes)
	osFile.Close()
	stat, err = os.Stat(file.Path)
	if err != nil {
		return true, fmt.Errorf("error reading file %s info: %w", file.Path, err)
	}
	file.changed = false
	file.time = stat.ModTime()
	return true, nil
}
