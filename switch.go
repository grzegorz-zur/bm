package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Switch is a mode for switching files.
type Switch struct {
	editor    *Editor
	query     Line
	paths     Paths
	selection Paths
	area      Area
	position  Position
}

// Show updates mode when switched to.
func (mode *Switch) Show() error {
	mode.query = Line{}
	var err error
	mode.paths, err = mode.read()
	if err != nil {
		return fmt.Errorf("error showing switch mode: %w", err)
	}
	sort.Sort(mode.paths)
	mode.filter()
	return nil
}

// Hide updates mode when switched from.
func (mode *Switch) Hide() error {
	return nil
}

// Key handles input events.
func (mode *Switch) Key(key Key) error {
	var err error
	switch key {
	case KeyTab:
		mode.editor.SwitchMode(mode.editor.Command)
	case KeyUp:
		mode.moveUp()
	case KeyDown:
		mode.moveDown()
	case KeyBackspace:
		mode.deletePreviousRune()
		mode.filter()
	case KeyEnter:
		err = mode.open()
		mode.editor.SwitchMode(mode.editor.Command)
	}
	if err != nil {
		return fmt.Errorf("error handling key %v: %w", key, err)
	}
	return nil
}

// Rune handles rune input.
func (mode *Switch) Rune(rune rune) error {
	mode.appendRune(rune)
	mode.filter()
	return nil
}

// Render renders mode.
func (mode *Switch) Render(content *Content) error {
	mode.area = mode.area.Resize(content.Size).Shift(mode.position)
	marked := len(mode.selection) > 0
	for l := mode.area.T; l < mode.area.B; l++ {
		rl := l - mode.area.T
		for c := mode.area.L; c < mode.area.R; c++ {
			rc := c - mode.area.L
			if l < len(mode.selection) {
				f := []rune(mode.selection[l])
				if c < len(f) {
					content.Runes[rl][rc] = f[c]
				}
			}
			if marked && l == mode.position.L {
				content.Marks[rl][rc] = true
			}
		}
	}
	status, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}
	content.Color = ColorBlue
	content.Position = mode.position
	content.Status = status
	content.Prompt = string(mode.query)
	content.Cursor = CursorPrompt
	return nil
}

func (mode *Switch) filter() {
	query := mode.query.String()
	mode.selection = make([]string, 0, len(mode.paths))
	for _, path := range mode.paths {
		if mode.match(path, query) {
			mode.selection = append(mode.selection, path)
		}
	}
	mode.position = Position{}
	return
}

func (mode *Switch) open() error {
	pos := mode.position
	path := mode.query.String()
	if pos.L < len(mode.selection) {
		path = mode.selection[pos.L]
	}
	err := mode.editor.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", path, err)
	}
	return nil
}

func (mode *Switch) appendRune(rune rune) {
	mode.query = mode.query.AppendRune(rune)
}

func (mode *Switch) deletePreviousRune() {
	mode.query = mode.query.DeletePreviousRune(len(mode.query))
}

func (mode *Switch) moveUp() {
	if mode.position.L > 0 {
		mode.position.L--
	}
}

func (mode *Switch) moveDown() {
	if mode.position.L+1 < len(mode.selection) {
		mode.position.L++
	}
}

func (mode *Switch) read() (paths []string, err error) {
	work, err := os.Getwd()
	if err != nil {
		return paths, fmt.Errorf("error reading working directory: %w", err)
	}
	walker := func(path string, info os.FileInfo, err error) error {
		relpath, err := filepath.Rel(work, path)
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			paths = append(paths, relpath)
		}
		return nil
	}
	err = filepath.Walk(work, walker)
	if err != nil {
		return paths, fmt.Errorf("error walking directory %s: %w", work, err)
	}
	return paths, nil
}

func (mode *Switch) match(path, query string) bool {
	if len(query) == 0 {
		return true
	}
	j := 0
	runes := []rune(query)
	for _, p := range path {
		q := runes[j]
		if p == q {
			j++
		}
		if j == len(query) {
			return true
		}
	}
	return false
}
