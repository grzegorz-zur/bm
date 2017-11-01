package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Input struct {
	*Editor
}

func (input *Input) Key(event tb.Event) (err error) {
	if event.Ch != 0 {
		input.Insert(event.Ch)
		return
	}
	switch event.Key {
	case tb.KeyCtrlQ: // TODO remove
		input.Quit()
	case tb.KeyArrowLeft:
		input.MoveLeft()
	case tb.KeyArrowRight:
		input.MoveRight()
	case tb.KeyArrowUp:
		input.MoveUp()
	case tb.KeyArrowDown:
		input.MoveDown()
	case tb.KeyDelete:
		input.Delete()
	}
	return
}
