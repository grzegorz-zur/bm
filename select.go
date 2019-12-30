package main

import (
	"fmt"
)

// Select mode.
type Select struct {
	editor *Editor
}

// Show updates mode when switched to.
func (m *Select) Show() error {
	m.editor.File.Select()
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
		m.editor.Motion(File.Left)
	case KeyRight:
		m.editor.Motion(File.Right)
	case KeyUp:
		m.editor.Motion(File.Up)
	case KeyDown:
		m.editor.Motion(File.Down)
	case KeyPageUp:
		m.editor.Motion(Paragraph(Backward))
	case KeyPageDown:
		m.editor.Motion(Paragraph(Forward))
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
		m.editor.Motion(File.Left)
	case 'f':
		m.editor.Motion(File.Right)
	case 'a':
		m.editor.Motion(File.Up)
	case 's':
		m.editor.Motion(File.Down)
	case 'e':
		m.editor.Motion(Word(Backward))
	case 'r':
		m.editor.Motion(Word(Forward))
	case 'q':
		m.editor.Motion(Paragraph(Backward))
	case 'w':
		m.editor.Motion(Paragraph(Forward))
	case 'g':
		m.editor.Copy()
		m.editor.SwitchMode(m.editor.Command)
	case 'j':
		m.editor.Copy()
		m.editor.Change(File.Delete)
		m.editor.SwitchMode(m.editor.Command)
	}
	if err != nil {
		return fmt.Errorf("error handling rune %v: %w", r, err)
	}
	return nil
}

// Render renders select mode.
func (m *Select) Render(cnt *Content) error {
	m.editor.File.Render(cnt, true)
	cnt.Color = ColorYellow
	cnt.Prompt = ""
	return nil
}
