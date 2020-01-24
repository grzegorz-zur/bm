package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Switch is a mode for switching files.
type Switch struct {
	editor   *Editor
	query    Line
	paths    Paths
	filtered Paths
	area     Area
	pos      Position
}

// Show updates mode when switched to.
func (m *Switch) Show() error {
	m.query = Line{}
	var err error
	m.paths, err = m.read()
	if err != nil {
		return fmt.Errorf("error showing switch mode: %w", err)
	}
	m.filter()
	return nil
}

// Hide updates mode when switched from.
func (m *Switch) Hide() error {
	return nil
}

// Key handles input events.
func (m *Switch) Key(k Key) error {
	var err error
	switch k {
	case KeyTab:
		if m.editor.Files.Empty() {
			m.editor.Quit()
		} else {
			m.editor.SwitchMode(m.editor.Command)
		}
	case KeyUp:
		m.moveUp()
	case KeyDown:
		m.moveDown()
	case KeyBackspace:
		m.deletePreviousRune()
		m.filter()
	case KeyEnter:
		err = m.open()
		m.editor.SwitchMode(m.editor.Command)
	}

	if err != nil {
		return fmt.Errorf("error handling key %v: %w", k, err)
	}

	return nil
}

// Rune handles rune input.
func (m *Switch) Rune(r rune) error {
	m.appendRune(r)
	m.filter()
	return nil
}

// Render renders mode.
func (m *Switch) Render(cnt *Content) error {
	m.area = m.area.Resize(cnt.Size).Shift(m.pos)
	marked := len(m.filtered) > 0
	for l := m.area.T; l < m.area.B; l++ {
		rl := l - m.area.T
		for c := m.area.L; c < m.area.R; c++ {
			rc := c - m.area.L
			if l < len(m.filtered) {
				f := []rune(m.filtered[l])
				if c < len(f) {
					cnt.Runes[rl][rc] = f[c]
				}
			}
			if marked && l == m.pos.L {
				cnt.Marks[rl][rc] = true
			}
		}
	}
	status, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}
	cnt.Color = ColorBlue
	cnt.Position = m.pos
	cnt.Status = status
	cnt.Prompt = string(m.query)
	cnt.Cursor = CursorPrompt
	return nil
}

func (m *Switch) filter() {
	query := m.query.String()
	m.filtered = make([]string, 0, len(m.paths))
	for _, path := range m.paths {
		if m.match(path, query) {
			m.filtered = append(m.filtered, path)
		}
	}
	sort.Sort(m.filtered)
	m.pos = Position{}
	return
}

func (m *Switch) open() error {
	p := m.pos
	path := m.query.String()
	if p.L < len(m.filtered) {
		path = m.filtered[p.L]
	}
	err := m.editor.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", path, err)
	}
	return nil
}

func (m *Switch) appendRune(r rune) {
	m.query = m.query.AppendRune(r)
}

func (m *Switch) deletePreviousRune() {
	m.query = m.query.DeletePreviousRune(len(m.query))
}

func (m *Switch) moveUp() {
	if m.pos.L > 0 {
		m.pos.L--
	}
}

func (m *Switch) moveDown() {
	if m.pos.L+1 < len(m.filtered) {
		m.pos.L++
	}
}

func (m *Switch) read() ([]string, error) {
	paths := []string{}
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

func (m *Switch) match(path, query string) bool {
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
