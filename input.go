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
		return
	}
	switch event.Key {
	case tb.KeyEsc:
		mode.Switch(mode.Editor.Normal)
	case tb.KeyArrowLeft:
		mode.MoveLeft()
	case tb.KeyArrowRight:
		mode.MoveRight()
	case tb.KeyArrowUp:
		mode.MoveUp()
	case tb.KeyArrowDown:
		mode.MoveDown()
	case tb.KeyEnter:
		mode.SplitLine()
	case tb.KeyDelete:
		mode.Delete()
	}
	return
}
