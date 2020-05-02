package main

import (
	"fmt"
)

// Input is a mode for typing.
type Input struct {
	editor *Editor
}

// Show updates mode when switched to.
func (mode *Input) Show() error {
	return nil
}

// Hide updates mode when switched from.
func (mode *Input) Hide() error {
	return nil
}

// Key handles input events.
func (mode *Input) Key(key Key) error {
	switch key {
	case KeyLeft:
		mode.editor.MoveLeft()
	case KeyRight:
		mode.editor.MoveRight()
	case KeyUp:
		mode.editor.MoveUp()
	case KeyDown:
		mode.editor.MoveDown()
	case KeyHome:
		mode.editor.MoveLineStart()
	case KeyEnd:
		mode.editor.MoveLineEnd()
	case KeyPageUp:
		mode.editor.MoveParagraphPrevious()
	case KeyPageDown:
		mode.editor.MoveParagraphNext()
	case KeyTab:
		mode.editor.Insert(string(Tab))
	case KeyEnter:
		mode.editor.Insert(string(EOL))
	case KeyBackspace:
		mode.editor.Backspace()
	case KeyDelete:
		mode.editor.Delete()
	case KeyCtrlSpace:
		mode.editor.SwitchMode(mode.editor.Command)
	}
	return nil
}

// Rune handles rune input.
func (mode *Input) Rune(rune rune) error {
	mode.editor.Insert(string(rune))
	return nil
}

// Render renders mode to the screen.
func (mode *Input) Render(view *View) error {
	err := mode.editor.File.Render(view)
	view.Color = ColorRed
	return fmt.Errorf("error rendering input mode: %w", err)
}
