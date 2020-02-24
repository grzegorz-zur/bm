package main

import (
	"fmt"
)

// Command mode.
type Command struct {
	editor *Editor
}

// Show updates mode when switched to.
func (mode *Command) Show() error {
	return nil
}

// Hide updates mode when switched from.
func (mode *Command) Hide() error {
	return nil
}

// Key handles special key.
func (mode *Command) Key(key Key) error {
	var err error
	switch key {
	case KeyTab:
		mode.editor.SwitchMode(mode.editor.Switch)
	}
	if mode.editor.Empty() {
		return nil
	}
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
func (mode *Command) Rune(rune rune) (err error) {
	switch rune {
	case 'b':
		err = mode.editor.Write()
		mode.editor.Pause()
	case 'B':
		err = mode.editor.Write()
		mode.editor.Quit()
	}
	if mode.editor.Empty() {
		return nil
	}
	switch rune {
	case ' ':
		mode.editor.SwitchMode(mode.editor.Input)
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
	case 'e':
		mode.editor.MoveWordPrevious()
	case 'r':
		mode.editor.MoveWordNext()
	case 'q':
		mode.editor.MoveParagraphPrevious()
	case 'w':
		mode.editor.MoveParagraphNext()
	case 'z':
		mode.editor.SwitchVersion(Backward)
	case 'x':
		mode.editor.SwitchVersion(Forward)
	case 'c':
		mode.editor.SwitchFile(Backward)
	case 'v':
		mode.editor.SwitchFile(Forward)
	case 'j':
		mode.editor.Delete()
	case 'J':
		mode.editor.DeleteLine()
	case 'k':
		mode.editor.LineBelow()
	case 'K':
		mode.editor.LineAbove()
	case 'g':
		mode.editor.SwitchMode(mode.editor.Select)
	case 'h':
		mode.editor.Paste()
	case 'n':
		_, err = mode.editor.File.Write()
		mode.editor.Files.Close()
	case 'N':
		_, err = mode.editor.Read(true)
	case 'm':
		err = mode.editor.Write()
	case 'M':
		_, err = mode.editor.File.Write()
	}
	if err != nil {
		return fmt.Errorf("error handling rune %v: %w", rune, err)
	}
	return nil
}

// Render renders mode.
func (mode *Command) Render(view *View) (err error) {
	if !mode.editor.Empty() {
		err = mode.editor.File.Render(view, false)
	}
	view.Color = ColorGreen
	view.Prompt = ""
	return err
}
