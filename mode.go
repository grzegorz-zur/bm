package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

type Mode interface {
	Visual
	Show() error
	Hide() error
	Key(event tb.Event) error
}

func renderNameAndPosition(name string, pos Position, color tb.Attribute, display *Display, bounds Bounds) (cursor Position, err error) {
	n := []rune(name)
	p := []rune(fmt.Sprintf("%d:%d", pos.Line+1, pos.Col+1))
	l := bounds.Right - bounds.Left
	for c := bounds.Left; c <= bounds.Right; c++ {
		r := ' '
		i := c - bounds.Left
		if i < len(n) {
			r = n[i]
		}
		j := i - l + len(p) - 1
		if 0 <= j && j < len(p) {
			r = p[j]
		}
		display.SetCell(c, bounds.Top, r, tb.ColorDefault|tb.AttrBold, color)
	}
	return
}
