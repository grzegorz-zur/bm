package bm

import (
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
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
		mode.Move(File.Left)
	case tb.KeyArrowRight:
		mode.Move(File.Right)
	case tb.KeyArrowUp:
		mode.Move(File.Up)
	case tb.KeyArrowDown:
		mode.Move(File.Down)
	}

	if err != nil {
		err = errors.Wrapf(err, "error handling event: %v", event)
	}

	return
}

func (mode *Command) Render(display *Display, bounds Bounds) (cursor Position, err error) {
	file, status := bounds.SplitHorizontal(-1)
	cursor, err = mode.File.Render(display, file)
	if err != nil {
		err = errors.Wrap(err, "error renderning file")
		return
	}
	_, err = mode.render(display, status)
	if err != nil {
		err = errors.Wrap(err, "error renderning status")
		return
	}
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
