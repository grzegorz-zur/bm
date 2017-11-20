package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Input struct {
	*Editor
}

func (mode *Input) Key(event tb.Event) (err error) {
	if event.Ch != 0 {
		mode.ApplyFileOp(InsertRune(event.Ch))
	}

	switch event.Key {
	case tb.KeyEsc:
		mode.SwitchMode(mode.Editor.Normal)
	case tb.KeyArrowLeft:
		mode.ApplyMoveOp(File.Left)
	case tb.KeyArrowRight:
		mode.ApplyMoveOp(File.Right)
	case tb.KeyArrowUp:
		mode.ApplyMoveOp(File.Up)
	case tb.KeyArrowDown:
		mode.ApplyMoveOp(File.Down)
	case tb.KeyEnter:
		mode.ApplyFileOp(File.Split)
	case tb.KeyDelete:
		mode.ApplyFileOp(File.DeleteRune)
	}

	return
}
