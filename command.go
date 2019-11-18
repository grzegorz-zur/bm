package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
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

// Key handles input events.
func (m *Command) Key(e tb.Event) error {

	var err error
	switch e.Ch {
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

	switch e.Key {
	case tb.KeySpace:
		m.SwitchMode(m.Editor.Input)
	case tb.KeyTab:
		m.SwitchMode(m.Editor.Switch)
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		err = m.WriteAll()
		m.Pause()
	case tb.KeyDelete:
		err = m.WriteAll()
		m.Quit()
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
	}

	if err != nil {
		return fmt.Errorf("error handling event %v: %w", e, err)
	}

	return nil
}

// Render renders mode to the screen.
func (m *Command) Render(d *Display, a Area) (Position, error) {
	file, status := a.SplitHorizontal(-1)
	cursor, err := m.File.Render(d, file)
	if err != nil {
		err = fmt.Errorf("error renderning file: %w", err)
		return cursor, err
	}
	_, err = renderNameAndPosition(m.Path, m.Position, tb.ColorGreen, d, status)
	if err != nil {
		err = fmt.Errorf("error renderning status: %w", err)
		return cursor, err
	}
	return cursor, nil
}
