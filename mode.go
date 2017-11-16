package bm

import (
	tb "github.com/nsf/termbox-go"
)

type Mode interface {
	Key(event tb.Event) error
}
