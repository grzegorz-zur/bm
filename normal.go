package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Normal struct {
	*Editor
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
		mode.Switch(mode.Editor.Input)
	case tb.KeyArrowLeft:
		mode.Move(File.Left)
	case tb.KeyArrowRight:
		mode.Move(File.Right)
	case tb.KeyArrowUp:
		mode.Move(File.Up)
	case tb.KeyArrowDown:
		mode.Move(File.Down)
	case tb.KeyDelete:
		mode.Change(File.DeleteRune)
	case tb.KeyCtrlD:
		mode.Next(Backward)
	case tb.KeyCtrlF:
		mode.Next(Forward)
	case tb.KeyCtrlQ:
		mode.Quit()
	case tb.KeyCtrlW:
		mode.Close()
	case tb.KeyCtrlE:
		mode.Write()
	case tb.KeyCtrlZ:
		mode.Stop()
	}

	return
}

func (mode *Normal) Display(bounds Bounds) (cursor Position, err error) {
	f, s := bounds.SplitHorizontal(-1)
	fc, err := mode.Current().Display(f)
	if err != nil {
		return
	}
	_, err = mode.display(s)
	if err != nil {
		return
	}
	cursor = fc
	return
}

func (mode *Normal) display(bounds Bounds) (cursor Position, err error) {
	name := []rune(mode.Current().Path)
	for c := bounds.Left; c <= bounds.Right; c++ {
		i := c - bounds.Left
		r := ' '
		if i < len(name) {
			r = name[i]
		}
		tb.SetCell(c, bounds.Top, r,
			tb.ColorDefault|tb.AttrBold,
			tb.ColorGreen)
	}
	return
}
