package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Input struct {
	*Editor
}

func (input *Input) Key(event tb.Event) (err error) {
	if event.Ch != 0 {
		input.InsertAfter(event.Ch)
		return
	}
	switch event.Key {
	case tb.KeyCtrlQ:
		input.Quit()
	}
	return
}
