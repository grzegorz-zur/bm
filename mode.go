package bm

import (
	tb "github.com/nsf/termbox-go"
)

type ModeType int

const (
	NormalMode ModeType = iota
	InsertMode
)

type Mode interface {
	Key(event tb.Event) error
}
