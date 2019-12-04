package main

import (
	"github.com/gdamore/tcell"
)

type Key int

const (
	KeyEscape Key = iota
	KeyTab
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
)

var keymap = map[tcell.Key]Key{
	tcell.KeyEscape:     KeyEscape,
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
}
