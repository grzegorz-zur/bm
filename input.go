package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

type Input struct {
	*Editor
}

// Show updates mode when switched to.
func (m *Input) Show() error {
	return nil
}

// Hide updates mode when switched from.
func (m *Input) Hide() error {
	return nil
}

// Key handles input events.
func (m *Input) Key(e tb.Event) error {
	if e.Ch != 0 {
		m.Change(InsertRune(e.Ch))
	}

	switch e.Key {
	case tb.KeyEsc:
		m.SwitchMode(m.Editor.Command)
	case tb.KeyArrowLeft:
		m.Motion(File.Left)
	case tb.KeyArrowRight:
		m.Motion(File.Right)
	case tb.KeyArrowUp:
		m.Motion(File.Up)
	case tb.KeyArrowDown:
		m.Motion(File.Down)
	case tb.KeyPgup:
		m.Motion(Paragraph(Backward))
	case tb.KeyPgdn:
		m.Motion(Paragraph(Forward))
	case tb.KeySpace:
		m.Change(InsertRune(' '))
	case tb.KeyTab:
		m.Change(InsertRune('\t'))
	case tb.KeyEnter:
		m.Change(File.Split)
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		m.Change(File.DeletePreviousRune)
	case tb.KeyDelete:
		m.Change(File.DeleteRune)
	}

	return nil
}

// Render renders mode.
func (m *Input) Render(d *Display, a Area) (Position, error) {
	file, status := a.SplitHorizontal(-1)
	cursor, err := m.File.Render(d, file)
	if err != nil {
		return cursor, fmt.Errorf("error rendering file: %w", err)
	}
	_, err = renderNameAndPosition(m.Path, m.Position, tb.ColorRed, d, status)
	if err != nil {
		return cursor, fmt.Errorf("error rendering status: %w", err)
	}
	return cursor, nil
}
