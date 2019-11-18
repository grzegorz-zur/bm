package main

import (
	"bufio"
	"fmt"
	"os"
)

// Open opens a file.
func Open(path string) (File, error) {
	file := File{
		Path:    path,
		History: &History{},
	}
	err := file.Load()
	if err != nil {
		return file, fmt.Errorf("error opening file %s: %w", path, err)
	}
	file.Archive()
	return file, nil
}

// RelaodIfModifed checks if file was modified outside and relaods it.
func (f *File) ReloadIfModified() (bool, error) {
	if f == nil {
		return false, ErrNoFile
	}
	modified, err := f.Modified()
	if err != nil {
		return modified, fmt.Errorf("error checking modification date on %s: %w", f.Path, err)
	}
	if !modified {
		return modified, nil
	}
	err = f.Reload()
	if err != nil {
		return modified, fmt.Errorf("error reloading file %s: %w", f.Path, err)
	}
	return modified, nil
}

// Modified checks if file was modified.
func (f *File) Modified() (bool, error) {
	if f == nil {
		return false, ErrNoFile
	}
	stat, err := os.Stat(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking file %s: %w", f.Path, err)
	}
	return stat.ModTime() != f.Time, nil
}

// Reload reloads the file.
func (f *File) Reload() error {
	if f == nil {
		return ErrNoFile
	}
	err := f.Load()
	if err != nil {
		return fmt.Errorf("error reloading file: %s: %w", f.Path, err)
	}
	f.Archive()
	return nil
}

// Load loads thie file.
func (f *File) Load() error {
	if f == nil {
		return ErrNoFile
	}
	flags := os.O_RDWR | os.O_CREATE
	fx, err := os.OpenFile(f.Path, flags, perm)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", f.Path, err)
	}
	defer fx.Close()
	err = f.update()
	if err != nil {
		return fmt.Errorf("error updating file information %s: %w", f.Path, err)
	}
	scanner := bufio.NewScanner(fx)
	f.Lines = nil
	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", f.Path, err)
		}
		line := scanner.Text()
		runes := []rune(line)
		f.Lines = append(f.Lines, runes)
	}
	return nil
}

// Write writes file contents.
func (f *File) Write() error {
	if f == nil {
		return ErrNoFile
	}
	fx, err := os.Create(f.Path)
	defer fx.Close()
	defer f.update()
	if err != nil {
		return fmt.Errorf("error writing file %s: %w", f.Path, err)
	}
	for _, runes := range f.Lines {
		line := string(runes)
		bytes := []byte(line)
		fx.Write(bytes)
		fx.Write(eol)
	}
	return nil
}
