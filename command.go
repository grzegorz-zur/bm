package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

type Command struct {
	*Editor
}

func (mode *Command) Show() (err error) {
	return
}

func (mode *Command) Hide() (err error) {
	return
}

func (mode *Command) Key(event tb.Event) (err error) {

	switch event.Ch {
	case 'd':
		mode.Motion(File.Left)
	case 'f':
		mode.Motion(File.Right)
	case 'a':
		mode.Motion(File.Up)
	case 's':
		mode.Motion(File.Down)
	case 'e':
		mode.Motion(mode.Word(Backward))
	case 'r':
		mode.Motion(mode.Word(Forward))
	case 'q':
		mode.Motion(mode.Paragraph(Backward))
	case 'w':
		mode.Motion(mode.Paragraph(Forward))
	case 'z':
		mode.SwitchVersion(Backward)
	case 'x':
		mode.SwitchVersion(Forward)
	case 'c':
		mode.SwitchFile(Backward)
	case 'v':
		mode.SwitchFile(Forward)
	case 'j':
		mode.Change(File.DeleteRune)
	case 'J':
		mode.Change(File.DeleteLine)
	case 'n':
		err = mode.Write()
		mode.Files.Close()
	case 'N':
		err = mode.Reload()
	case 'm':
		err = mode.WriteAll()
	case 'M':
		err = mode.Write()
	}

	switch event.Key {
	case tb.KeySpace:
		mode.SwitchMode(mode.Editor.Input)
	case tb.KeyTab:
		mode.SwitchMode(mode.Editor.Switch)
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		err = mode.WriteAll()
		mode.Pause()
	case tb.KeyDelete:
		err = mode.WriteAll()
		mode.Quit()
	case tb.KeyArrowLeft:
		mode.Motion(File.Left)
	case tb.KeyArrowRight:
		mode.Motion(File.Right)
	case tb.KeyArrowUp:
		mode.Motion(File.Up)
	case tb.KeyArrowDown:
		mode.Motion(File.Down)
	case tb.KeyPgup:
		mode.Motion(mode.Paragraph(Backward))
	case tb.KeyPgdn:
		mode.Motion(mode.Paragraph(Forward))
	}

	if err != nil {
		err = fmt.Errorf("error handling event %v: %w", event, err)
	}

	return
}

func (mode *Command) Render(display *Display, bounds Bounds) (cursor Position, err error) {
	file, status := bounds.SplitHorizontal(-1)
	cursor, err = mode.File.Render(display, file)
	if err != nil {
		err = fmt.Errorf("error renderning file: %w", err)
		return
	}
	_, err = renderNameAndPosition(mode.Path, mode.Position, tb.ColorGreen, display, status)
	if err != nil {
		err = fmt.Errorf("error renderning status: %w", err)
		return
	}
	return
}
