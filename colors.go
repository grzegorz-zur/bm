package main

import (
	"github.com/gdamore/tcell"
)

// Color is a color used here.
type Color int

// Colors used in editor.
const (
	ColorNone Color = iota
	ColorRed
	ColorGreen
	ColorBlue
	ColorYellow
)

var colors = map[Color]tcell.Color{
	ColorRed:    tcell.ColorRed,
	ColorGreen:  tcell.ColorGreen,
	ColorBlue:   tcell.ColorBlue,
	ColorYellow: tcell.ColorYellow,
}
