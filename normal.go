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
		mode.MoveLeft()
	case 'f':
		mode.MoveRight()
	case 'k':
		mode.MoveUp()
	case 'j':
		mode.MoveDown()
	}
	switch event.Key {
	case tb.KeySpace:
		mode.Switch(mode.Editor.Input)
	case tb.KeyArrowLeft:
		mode.MoveLeft()
	case tb.KeyArrowRight:
		mode.MoveRight()
	case tb.KeyArrowUp:
		mode.MoveUp()
	case tb.KeyArrowDown:
		mode.MoveDown()
	case tb.KeyDelete:
		mode.Delete()
	case tb.KeyCtrlQ:
		mode.Quit()
	case tb.KeyCtrlW:
		mode.Write()
	}
	return
}
