package main

import (
	"fmt"
)

// Select mode.
type Select struct {
	editor *Editor
}

// Show updates mode when switched to.
func (mode *Select) Show() error {
	mode.editor.File.Select()
	return nil
}

// Hide updates mode when switched from.
func (mode *Select) Hide() error {
	return nil
}

// Key handles special key.
func (mode *Select) Key(key Key) error {
	var err error
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
	}
	if err != nil {
		return fmt.Errorf("error handling key %v: %w", key, err)
	}
	return nil
}

// Rune handles runes.
func (mode *Select) Rune(rune rune) error {
	var err error
	switch rune {
	case 'd':
		mode.editor.MoveLeft()
	case 'f':
		mode.editor.MoveRight()
	case 'a':
		mode.editor.MoveUp()
	case 's':
		mode.editor.MoveDown()
	case 'D':
		mode.editor.MoveLineStart()
	case 'F':
		mode.editor.MoveLineEnd()
	case 'A':
		mode.editor.MoveFileStart()
	case 'S':
		mode.editor.MoveFileEnd()
	case 'q':
		mode.editor.MoveParagraphPrevious()
	case 'w':
		mode.editor.MoveParagraphNext()
	case 'e':
		mode.editor.MoveWordPrevious()
	case 'r':
		mode.editor.MoveWordNext()
	case 'H':
		mode.editor.SwitchMode(mode.editor.Command)
	case 'h':
		mode.editor.Copy()
	case 'j':
		mode.editor.Cut()
	}
	if err != nil {
		return fmt.Errorf("error handling rune %v: %w", rune, err)
	}
	return nil
}

// Render renders select mode.
func (mode *Select) Render(view *View) error {
	view.Select = true
	err := mode.editor.File.Render(view)
	if err != nil {
		return fmt.Errorf("error rendering select mode: %w", err)
	}
	view.Color = ColorYellow
	view.Prompt = ""
	return nil
}
