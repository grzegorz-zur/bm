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
	case tb.KeyCtrlQ:
		mode.Quit()
	case tb.KeyCtrlW:
		mode.Write()
	}

	return
}
