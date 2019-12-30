package main

import (
	"github.com/gdamore/tcell"
)

// Key represends key of an keyboard.
type Key int

// Keys used in editor.
const (
	KeyTab Key = iota
	KeyEnter
	KeyBackspace
	KeyDelete
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyCtrlSpace
)

var keymap = map[tcell.Key]Key{
	tcell.KeyTab:        KeyTab,
	tcell.KeyEnter:      KeyEnter,
	tcell.KeyBackspace:  KeyBackspace,
	tcell.KeyBackspace2: KeyBackspace,
	tcell.KeyDelete:     KeyDelete,
	tcell.KeyUp:         KeyUp,
	tcell.KeyDown:       KeyDown,
	tcell.KeyLeft:       KeyLeft,
	tcell.KeyRight:      KeyRight,
	tcell.KeyHome:       KeyHome,
	tcell.KeyEnd:        KeyEnd,
	tcell.KeyPgUp:       KeyPageUp,
	tcell.KeyPgDn:       KeyPageDown,
	tcell.KeyCtrlSpace:  KeyCtrlSpace,
}
