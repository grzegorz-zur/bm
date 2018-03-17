package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Normal struct {
	*Editor
}

func (mode *Normal) Show() (err error) {
	return
}

func (mode *Normal) Hide() (err error) {
	return
}

func (mode *Normal) Key(event tb.Event) (err error) {
	switch event.Ch {
	case 'd':
		mode.Move(File.Left)
	case 'f':
		mode.Move(File.Right)
	case 'k':
		mode.Move(File.Up)
	case 'j':
		mode.Move(File.Down)
	}

	switch event.Key {
	case tb.KeySpace:
		mode.SwitchMode(mode.Editor.Input)
	case tb.KeyTab:
		mode.SwitchMode(mode.Editor.Switch)
	case tb.KeyArrowLeft:
		mode.Move(File.Left)
	case tb.KeyArrowRight:
		mode.Move(File.Right)
	case tb.KeyArrowUp:
		mode.Move(File.Up)
	case tb.KeyArrowDown:
		mode.Move(File.Down)
	case tb.KeyBackspace:
	case tb.KeyBackspace2:
		mode.Change(File.DeletePreviousRune)
	case tb.KeyDelete:
		mode.Change(File.DeleteRune)
	case tb.KeyCtrlD:
		mode.Next(Backward)
	case tb.KeyCtrlF:
		mode.Next(Forward)
	case tb.KeyCtrlQ:
		mode.Quit()
	case tb.KeyCtrlW:
		mode.Files.Close()
	case tb.KeyCtrlE:
		mode.Write()
	case tb.KeyCtrlZ:
		mode.Stop()
	}

	return
}

func (mode *Normal) Render(display *Display, bounds Bounds) (cursor Position, err error) {
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

func (mode *Normal) render(display *Display, bounds Bounds) (cursor Position, err error) {
	if mode.File == nil {
		return
	}
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
