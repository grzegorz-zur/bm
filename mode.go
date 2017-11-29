package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Mode interface {
	Visual
	Key(event tb.Event) error
}
