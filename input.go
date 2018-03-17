package bm

import (
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
		mode.SwitchMode(mode.Editor.Normal)
	case tb.KeyArrowLeft:
		mode.Move(File.Left)
	case tb.KeyArrowRight:
		mode.Move(File.Right)
	case tb.KeyArrowUp:
		mode.Move(File.Up)
	case tb.KeyArrowDown:
		mode.Move(File.Down)
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
	f, s := bounds.SplitHorizontal(-1)
	fc, err := mode.File.Render(display, f)
	if err != nil {
		return
	}
	_, err = mode.render(display, s)
	if err != nil {
		return
	}
	cursor = fc
	return
}

func (mode *Input) render(display *Display, bounds Bounds) (cursor Position, err error) {
	name := []rune(mode.Path)
	for c := bounds.Left; c <= bounds.Right; c++ {
		i := c - bounds.Left
		r := ' '
		if i < len(name) {
			r = name[i]
		}
		display.SetCell(c, bounds.Top, r, tb.ColorDefault|tb.AttrBold, tb.ColorRed)
	}
	return
}
