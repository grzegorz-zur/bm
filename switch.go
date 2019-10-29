package bm

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
	"os"
	"path/filepath"
	"strings"
)

type Switch struct {
	*Editor
	query    Line
	paths    []string
	filtered []string
	position Position
	window   Bounds
}

func (mode *Switch) Show() (err error) {
	mode.query = Line{}
	mode.paths, err = mode.read()
	if err != nil {
		err = fmt.Errorf("error showing switch mode: %w", err)
		return
	}
	mode.filter()
	return
}

func (mode *Switch) Hide() (err error) {
	return
}

func (mode *Switch) Key(event tb.Event) (err error) {
	if event.Ch != 0 {
		mode.appendRune(event.Ch)
		mode.filter()
	}

	switch event.Key {
	case tb.KeyEsc:
		if mode.Files.Empty() {
			mode.Quit()
		} else {
			mode.SwitchMode(mode.Command)
		}
	case tb.KeyArrowUp:
		mode.moveUp()
	case tb.KeyArrowDown:
		mode.moveDown()
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		mode.deletePreviousRune()
		mode.filter()
	case tb.KeyEnter:
		err = mode.open()
		mode.SwitchMode(mode.Command)
	}

	if err != nil {
		err = fmt.Errorf("error handling event %v: %w", event, err)
	}

	return
}

func (mode *Switch) filter() {
	query := mode.query.String()
	mode.filtered = make([]string, 0, len(mode.paths))
	for _, path := range mode.paths {
		if match(path, query) {
			mode.filtered = append(mode.filtered, path)
		}
	}
	mode.position = Position{}
	return
}

func (mode *Switch) open() (err error) {
	pos := mode.position
	path := mode.query.String()
	if pos.Line < len(mode.filtered) {
		path = mode.filtered[pos.Line]
	}
	err = mode.Open(path)
	if err != nil {
		err = fmt.Errorf("error opening file %s: %w", path, err)
	}
	return
}

func (mode *Switch) appendRune(r rune) {
	mode.query = mode.query.AppendRune(r)
}

func (mode *Switch) deletePreviousRune() {
	mode.query = mode.query.DeletePreviousRune(len(mode.query))
}

func (mode *Switch) moveUp() {
	p := mode.position
	if p.Line > 0 {
		mode.position = Position{Line: p.Line - 1}
	}
}

func (mode *Switch) moveDown() {
	p := mode.position
	if p.Line+1 < len(mode.filtered) {
		mode.position = Position{Line: p.Line + 1}
	}
}

func (mode *Switch) Render(display *Display, bounds Bounds) (cursor Position, err error) {
	paths, status := bounds.SplitHorizontal(-1)
	err = mode.renderPaths(display, paths)
	if err != nil {
		err = fmt.Errorf("error rendering paths: %w", err)
		return
	}
	cursor, err = mode.renderInput(display, status)
	if err != nil {
		err = fmt.Errorf("error rendering status: %w", err)
		return
	}
	return
}

func (mode *Switch) renderPaths(display *Display, bounds Bounds) (err error) {
	paths := mode.filtered
	mode.scroll()
	size := bounds.Size()
	mode.size(size)
	p := mode.position
	w := mode.window
	for line := w.Top; line <= w.Bottom; line++ {
		if line >= len(paths) {
			break
		}
		foreground := tb.ColorDefault
		background := tb.ColorDefault
		if line == p.Line {
			foreground = tb.ColorBlack
			background = tb.ColorWhite
		}
		path := paths[line]
		runes := []rune(path)
		screenLine := bounds.Top + line - w.Top
		for col := w.Left; col <= w.Right; col++ {
			if col >= len(runes) {
				break
			}
			symbol := runes[col]
			screenCol := bounds.Left + col - w.Left
			display.SetCell(screenCol, screenLine, symbol, foreground, background)
		}
	}
	return
}

func (mode *Switch) size(size Size) {
	w := &mode.window
	w.Bottom = w.Top + size.Lines
	w.Right = w.Left + size.Cols
	return
}

func (mode *Switch) scroll() {
	p := mode.position
	w := &mode.window
	height := w.Bottom - w.Top
	width := w.Right - w.Left

	switch {
	case p.Line < w.Top:
		w.Top = p.Line
		w.Bottom = w.Top + height
	case p.Line > w.Bottom:
		w.Bottom = p.Line
		w.Top = w.Bottom - height
	}

	switch {
	case p.Col < w.Left:
		w.Left = p.Col
		w.Right = w.Left + width
	case p.Col > w.Right:
		w.Right = p.Col
		w.Left = w.Right - width
	}
}

func (mode *Switch) renderInput(display *Display, bounds Bounds) (cursor Position, err error) {
	for c := bounds.Left; c <= bounds.Right; c++ {
		i := c - bounds.Left
		r := ' '
		if i < len(mode.query) {
			r = mode.query[i]
		}
		display.SetCell(c, bounds.Top, r, tb.ColorDefault|tb.AttrBold, tb.ColorBlue)
	}
	cursor = Position{Line: bounds.Top, Col: len(mode.query)}
	return
}

func (mode *Switch) read() (paths []string, err error) {
	work, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("error reading working directory: %w", err)
		return
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
		err = fmt.Errorf("error walking directory %s: %w", work, err)
	}
	return
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
