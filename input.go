package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Input struct {
	*Editor
}

func (mode *Input) Key(event tb.Event) (err error) {
	if event.Ch != 0 {
		mode.Insert(event.Ch)
	}

	switch event.Key {
	case tb.KeyEsc:
		mode.Switch(mode.Editor.Normal)
	case tb.KeyArrowLeft:
		mode.ApplyFileOp(MoveOp(Left))
	case tb.KeyArrowRight:
		mode.ApplyFileOp(MoveOp(Right))
	case tb.KeyArrowUp:
		mode.ApplyFileOp(MoveOp(Up))
	case tb.KeyArrowDown:
		mode.ApplyFileOp(MoveOp(Down))
	case tb.KeyEnter:
		mode.ApplyFileOp(Split)
	case tb.KeyDelete:
		mode.Delete()
	}

	return
}
