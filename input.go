package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

type Input struct {
	*Editor
}

func (mode *Input) Show() (err error) {
	return
}

func (mode *Input) Hide() (err error) {
	return
}

func (mode *Input) Key(event tb.Event) (err error) {
	if event.Ch != 0 {
		mode.Change(InsertRune(event.Ch))
	}

	switch event.Key {
	case tb.KeyEsc:
		mode.SwitchMode(mode.Editor.Command)
	case tb.KeyArrowLeft:
		mode.Motion(File.Left)
	case tb.KeyArrowRight:
		mode.Motion(File.Right)
	case tb.KeyArrowUp:
		mode.Motion(File.Up)
	case tb.KeyArrowDown:
		mode.Motion(File.Down)
	case tb.KeySpace:
		mode.Change(InsertRune(' '))
	case tb.KeyTab:
		mode.Change(InsertRune('\t'))
	case tb.KeyEnter:
		mode.Change(File.Split)
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		mode.Change(File.DeletePreviousRune)
	case tb.KeyDelete:
		mode.Change(File.DeleteRune)
	}

	return
}

func (mode *Input) Render(display *Display, bounds Bounds) (cursor Position, err error) {
	file, status := bounds.SplitHorizontal(-1)
	cursor, err = mode.File.Render(display, file)
	if err != nil {
		err = fmt.Errorf("error rendering file: %w", err)
		return
	}
	_, err = renderNameAndPosition(mode.Path, mode.Position, tb.ColorRed, display, status)
	if err != nil {
		err = fmt.Errorf("error rendering status: %w", err)
		return
	}
	return
}
