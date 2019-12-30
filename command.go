package main

import (
	"fmt"
)

// Command mode.
type Command struct {
	editor *Editor
}

// Show updates mode when switched to.
func (m *Command) Show() error {
	return nil
}

// Hide updates mode when switched from.
func (m *Command) Hide() error {
	return nil
}

// Key handles special key.
func (m *Command) Key(k Key) error {
	var err error
	switch k {
	case KeyTab:
		m.editor.SwitchMode(m.editor.Switch)
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
func (m *Command) Rune(r rune) error {
	var err error
	switch r {
	case ' ':
		m.editor.SwitchMode(m.editor.Input)
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
	case 'z':
		m.editor.SwitchVersion(Backward)
	case 'x':
		m.editor.SwitchVersion(Forward)
	case 'c':
		m.editor.SwitchFile(Backward)
	case 'v':
		m.editor.SwitchFile(Forward)
	case 'j':
		m.editor.Change(File.DeleteRune)
	case 'J':
		m.editor.Change(File.DeleteLine)
	case 'g':
		m.editor.SwitchMode(m.editor.Select)
	case 'h':
		m.editor.PasteBlock()
	case 'H':
		m.editor.PasteInline()
	case 'n':
		err = m.editor.Write()
		m.editor.Files.Close()
	case 'N':
		err = m.editor.Reload()
	case 'm':
		err = m.editor.WriteAll()
	case 'M':
		err = m.editor.Write()
	case 'b':
		err = m.editor.WriteAll()
		m.editor.Pause()
	case 'B':
		err = m.editor.WriteAll()
		m.editor.Quit()
	}
	if err != nil {
		return fmt.Errorf("error handling rune %v: %w", r, err)
	}
	return nil
}

// Render renders mode.
func (m *Command) Render(cnt *Content) error {
	m.editor.File.Render(cnt, false)
	cnt.Color = ColorGreen
	cnt.Prompt = ""
	return nil
}
