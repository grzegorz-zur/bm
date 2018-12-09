package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Display struct{}

func (display *Display) Clear(fg, bg tb.Attribute) error {
	if display == nil {
		return nil
	}
	return tb.Clear(fg, bg)
}

func (display *Display) Close() {
	if display == nil {
		return
	}
	tb.Close()
}

func (display *Display) Flush() error {
	if display == nil {
		return nil
	}
	return tb.Flush()
}

func (display *Display) Init() error {
	if display == nil {
		return nil
	}
	return tb.Init()
}

func (display *Display) SetCell(x, y int, ch rune, fg, bg tb.Attribute) {
	if display == nil {
		return
	}
	tb.SetCell(x, y, ch, fg, bg)
}

func (display *Display) SetCursor(x, y int) {
	if display == nil {
		return
	}
	tb.SetCursor(x, y)
}

func (display *Display) Size() (width, height int) {
	if display == nil {
		return
	}
	return tb.Size()
}