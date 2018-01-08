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
		mode.ApplyMoveOp(File.Left)
	case 'f':
		mode.ApplyMoveOp(File.Right)
	case 'k':
		mode.ApplyMoveOp(File.Up)
	case 'j':
		mode.ApplyMoveOp(File.Down)
	}

	switch event.Key {
	case tb.KeySpace:
		mode.SwitchMode(mode.Editor.Input)
	case tb.KeyArrowLeft:
		mode.ApplyMoveOp(File.Left)
	case tb.KeyArrowRight:
		mode.ApplyMoveOp(File.Right)
	case tb.KeyArrowUp:
		mode.ApplyMoveOp(File.Up)
	case tb.KeyArrowDown:
		mode.ApplyMoveOp(File.Down)
	case tb.KeyDelete:
		mode.ApplyFileOp(File.DeleteRune)
	case tb.KeyCtrlD:
		mode.SwitchFile(mode.Next(mode.File, Backward))
	case tb.KeyCtrlF:
		mode.SwitchFile(mode.Next(mode.File, Forward))
	case tb.KeyCtrlQ:
		mode.Quit()
	case tb.KeyCtrlW:
		mode.Write()
	case tb.KeyCtrlZ:
		mode.Stop()
	}

	return
}

func (mode *Normal) Display(bounds Bounds) (cursor Position, err error) {
	f, s := bounds.SplitHorizontal(-1)
	fc, err := mode.Editor.File.Display(f)
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
	name := []rune(mode.Editor.File.Path)
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
