package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
	"os"
	"path/filepath"
	"strings"
)

// Switch is a mode for switching files.
type Switch struct {
	*Editor
	query    Line
	paths    []string
	filtered []string
	position Position
	window   Area
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
func (m *Switch) Key(e tb.Event) error {
	if e.Ch != 0 {
		m.appendRune(e.Ch)
		m.filter()
	}
	var err error
	switch e.Key {
	case tb.KeyEsc:
		if m.Files.Empty() {
			m.Quit()
		} else {
			m.SwitchMode(m.Command)
		}
	case tb.KeyArrowUp:
		m.moveUp()
	case tb.KeyArrowDown:
		m.moveDown()
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		m.deletePreviousRune()
		m.filter()
	case tb.KeyEnter:
		err = m.open()
		m.SwitchMode(m.Command)
	}

	if err != nil {
		return fmt.Errorf("error handling event %v: %w", e, err)
	}

	return nil
}

func (m *Switch) filter() {
	query := m.query.String()
	m.filtered = make([]string, 0, len(m.paths))
	for _, path := range m.paths {
		if match(path, query) {
			m.filtered = append(m.filtered, path)
		}
	}
	return
}

func (m *Switch) open() error {
	p := m.position
	path := m.query.String()
	if p.L < len(m.filtered) {
		path = m.filtered[p.L]
	}
	err := m.Open(path)
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
	if m.position.L > 0 {
		m.position.L--
	}
}

func (m *Switch) moveDown() {
	if m.position.L+1 < len(m.filtered) {
		m.position.L++
	}
}

// Render renders list of files.
func (m *Switch) Render(d *Display, a Area) (Position, error) {
	paths, status := a.SplitHorizontal(-1)
	err := m.renderPaths(d, paths)
	if err != nil {
		return Position{}, fmt.Errorf("error rendering paths: %w", err)
	}
	cursor, err := m.renderInput(d, status)
	if err != nil {
		return cursor, fmt.Errorf("error rendering status: %w", err)
	}
	return cursor, nil
}

func (m *Switch) renderPaths(d *Display, a Area) error {
	paths := m.filtered
	m.scroll()
	s := a.Size()
	m.size(s)
	p := m.position
	w := m.window
	for l := w.T; l <= w.B; l++ {
		if l >= len(paths) {
			break
		}
		fg := tb.ColorDefault
		bg := tb.ColorDefault
		if l == p.L {
			fg = tb.ColorBlack
			bg = tb.ColorWhite
		}
		path := paths[l]
		runes := []rune(path)
		sl := a.T + l - w.T
		for c := w.L; c <= w.R; c++ {
			if c >= len(runes) {
				break
			}
			r := runes[c]
			sc := a.L + c - w.L
			d.SetCell(sc, sl, r, fg, bg)
		}
	}
	return nil
}

func (m *Switch) size(s Size) {
	w := &m.window
	w.B = w.T + s.L - 1
	w.R = w.L + s.C - 1
}

func (m *Switch) scroll() {
	p := m.position
	w := &m.window
	s := w.Size()

	switch {
	case p.L < w.T:
		w.T = p.L
		w.B = w.T + s.L - 1
	case p.L > w.B:
		w.B = p.L
		w.T = w.B - s.L + 1
	}

	switch {
	case p.C < w.L:
		w.L = p.C
		w.R = w.L + s.C - 1
	case p.C > w.R:
		w.R = p.C
		w.L = w.R - s.C + 1
	}
}

func (m *Switch) renderInput(d *Display, a Area) (Position, error) {
	for c := a.L; c <= a.R; c++ {
		i := c - a.L
		r := ' '
		if i < len(m.query) {
			r = m.query[i]
		}
		d.SetCell(c, a.T, r, tb.ColorDefault|tb.AttrBold, tb.ColorBlue)
	}
	return Position{L: a.T, C: len(m.query)}, nil
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
		if include(relpath, info) {
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

func include(path string, info os.FileInfo) bool {
	if strings.HasPrefix(path, ".") || strings.Contains(path, "/.") {
		return false
	}
	if !info.Mode().IsRegular() {
		return false
	}
	return true
}

func match(path, query string) bool {
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
