package main

import (
	tb "github.com/nsf/termbox-go"
)

// Display is a terminal.
//
// TODO rewrite to seperate rest of the code from display library.
type Display struct{}

func (d *Display) Clear(fg, bg tb.Attribute) error {
	if d == nil {
		return nil
	}
	return tb.Clear(fg, bg)
}

func (d *Display) Close() {
	if d == nil {
		return
	}
	tb.Close()
}

func (d *Display) Flush() error {
	if d == nil {
		return nil
	}
	return tb.Flush()
}

func (d *Display) Init() error {
	if d == nil {
		return nil
	}
	return tb.Init()
}

func (d *Display) SetCell(x, y int, ch rune, fg, bg tb.Attribute) {
	if d == nil {
		return
	}
	tb.SetCell(x, y, ch, fg, bg)
}

func (d *Display) SetCursor(x, y int) {
	if d == nil {
		return
	}
	tb.SetCursor(x, y)
}

func (d *Display) Size() (x, y int) {
	if d == nil {
		return
	}
	return tb.Size()
}
