package main

import (
	"fmt"
)

// Command mode.
type Command struct {
	*Editor
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
		m.SwitchMode(m.Editor.Switch)
	case KeyBackspace:
		err = m.WriteAll()
		m.Pause()
	case KeyDelete:
		err = m.WriteAll()
		m.Quit()
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
func (m *Command) Rune(r rune) error {
	var err error
	switch r {
	case ' ':
		m.SwitchMode(m.Editor.Input)
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
	case 'z':
		m.SwitchVersion(Backward)
	case 'x':
		m.SwitchVersion(Forward)
	case 'c':
		m.SwitchFile(Backward)
	case 'v':
		m.SwitchFile(Forward)
	case 'j':
		m.Change(File.DeleteRune)
	case 'J':
		m.Change(File.DeleteLine)
	case 'n':
		err = m.Write()
		m.Files.Close()
	case 'N':
		err = m.Reload()
	case 'm':
		err = m.WriteAll()
	case 'M':
		err = m.Write()
	}
	if err != nil {
		return fmt.Errorf("error handling rune %v: %w", r, err)
	}
	return nil
}

func (m *Command) Render(cnt *Content) error {
	m.File.Render(cnt)
	cnt.Color = ColorGreen
	cnt.Prompt = ""
	return nil
}
