package main

import (
	tb "github.com/nsf/termbox-go"
)

// Mode represents a mode of the editor.
//
// Mode is responsible for input and display.
type Mode interface {
	Visual
	Show() error
	Hide() error
	Key(event tb.Event) error
}
