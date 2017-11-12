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
		mode.ApplyFileOp(MoveOp(Left))
	case 'f':
		mode.ApplyFileOp(MoveOp(Right))
	case 'k':
		mode.ApplyFileOp(MoveOp(Up))
	case 'j':
		mode.ApplyFileOp(MoveOp(Down))
	}

	switch event.Key {
	case tb.KeySpace:
		mode.Switch(mode.Editor.Input)
	case tb.KeyArrowLeft:
		mode.ApplyFileOp(MoveOp(Left))
	case tb.KeyArrowRight:
		mode.ApplyFileOp(MoveOp(Right))
	case tb.KeyArrowUp:
		mode.ApplyFileOp(MoveOp(Up))
	case tb.KeyArrowDown:
		mode.ApplyFileOp(MoveOp(Down))
	case tb.KeyDelete:
		mode.Delete()
	case tb.KeyCtrlQ:
		mode.Quit()
	case tb.KeyCtrlW:
		mode.Write()
	}

	return
}
