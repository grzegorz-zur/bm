package main

import (
	tb "github.com/nsf/termbox-go"
)

type Mode interface {
	Visual
	Show() error
	Hide() error
	Key(event tb.Event) error
}
