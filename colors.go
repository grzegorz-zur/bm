package main

import (
	"github.com/gdamore/tcell"
)

type Color int

const (
	ColorNone Color = iota
	ColorRed
	ColorGreen
	ColorBlue
)

var colors = map[Color]tcell.Color{
	ColorRed:   tcell.ColorRed,
	ColorGreen: tcell.ColorGreen,
	ColorBlue:  tcell.ColorBlue,
}
