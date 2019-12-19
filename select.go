package main

import (
	"fmt"
)

// Select mode.
type Select struct {
	*Editor
}

// Show updates mode when switched to.
func (m *Select) Show() error {
	m.File.Select()
	return nil
}

// Hide updates mode when switched from.
func (m *Select) Hide() error {
	return nil
}

// Key handles special key.
func (m *Select) Key(k Key) error {
	var err error
	switch k {
	case KeyLeft:
		m.Motion(File.Left)
	case KeyRight:
		m.Motion(File.Right)
	case KeyUp:
		m.Motion(File.Up)
	case KeyDown:
		m.Motion(File.Down)
	case KeyPageUp:
		m.Motion(Paragraph(Backward))
	case KeyPageDown:
		m.Motion(Paragraph(Forward))
	}
	if err != nil {
		return fmt.Errorf("error handling key %v: %w", k, err)
	}
	return nil
}

// Rune handles runes.
func (m *Select) Rune(r rune) error {
	var err error
	switch r {
	case 'd':
		m.Motion(File.Left)
	case 'f':
		m.Motion(File.Right)
	case 'a':
		m.Motion(File.Up)
	case 's':
		m.Motion(File.Down)
	case 'e':
		m.Motion(Word(Backward))
	case 'r':
		m.Motion(Word(Forward))
	case 'q':
		m.Motion(Paragraph(Backward))
	case 'w':
		m.Motion(Paragraph(Forward))
	case 'g':
		m.Copy()
		m.SwitchMode(m.Command)
	case 'j':
		m.Copy()
		m.Change(File.Delete)
		m.SwitchMode(m.Command)
	}
	if err != nil {
		return fmt.Errorf("error handling rune %v: %w", r, err)
	}
	return nil
}

// Render renders select mode.
func (m *Select) Render(cnt *Content) error {
	m.File.Render(cnt, true)
	cnt.Color = ColorYellow
	cnt.Prompt = ""
	return nil
}
