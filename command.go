package bm

import (
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
		mode.Move(File.Left)
	case 'f':
		mode.Move(File.Right)
	case 'a':
		mode.Move(File.Up)
	case 's':
		mode.Move(File.Down)
	case 'c':
		mode.Next(Backward)
	case 'v':
		mode.Next(Forward)
	case 'j':
		mode.Change(File.DeleteRune)
	case 'J':
		mode.Change(File.DeleteLine)
	case 'n':
		mode.WriteAll()
		mode.Files.Close()
	case 'm':
		mode.WriteAll()
	case 'M':
		mode.Write()
	}

	switch event.Key {
	case tb.KeySpace:
		mode.SwitchMode(mode.Editor.Input)
	case tb.KeyTab:
		mode.SwitchMode(mode.Editor.Switch)
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		mode.WriteAll()
		mode.Pause()
	case tb.KeyDelete:
		mode.WriteAll()
		mode.Quit()
	case tb.KeyArrowLeft:
		mode.Move(File.Left)
	case tb.KeyArrowRight:
		mode.Move(File.Right)
	case tb.KeyArrowUp:
		mode.Move(File.Up)
	case tb.KeyArrowDown:
		mode.Move(File.Down)
	}

	return
}

func (mode *Command) Render(display *Display, bounds Bounds) (cursor Position, err error) {
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

func (mode *Command) render(display *Display, bounds Bounds) (cursor Position, err error) {
	name := []rune(mode.Path)
	for c := bounds.Left; c <= bounds.Right; c++ {
		i := c - bounds.Left
		r := ' '
		if i < len(name) {
			r = name[i]
		}
		display.SetCell(c, bounds.Top, r, tb.ColorDefault|tb.AttrBold, tb.ColorGreen)
	}
	return
}
