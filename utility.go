package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

func wrap(pos, len, step int, d Direction) int {
	if len == 0 {
		return 0
	}
	return ((pos+step*d.Value())%len + len) % len
}

// TODO this should become a method of Display
func renderNameAndPosition(name string, pos Position, color tb.Attribute, display *Display, area Area) (Position, error) {
	n := []rune(name)
	p := []rune(fmt.Sprintf("%d:%d", pos.L+1, pos.C+1))
	l := area.R - area.L
	for c := area.L; c <= area.R; c++ {
		r := ' '
		i := c - area.L
		if i < len(n) {
			r = n[i]
		}
		j := i - l + len(p) - 1
		if 0 <= j && j < len(p) {
			r = p[j]
		}
		display.SetCell(c, area.T, r, tb.ColorDefault|tb.AttrBold, color)
	}
	return Position{area.T, area.L}, nil
}
