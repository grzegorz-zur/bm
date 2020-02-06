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
		mode.editor.Motion(File.Left)
	case KeyRight:
		mode.editor.Motion(File.Right)
	case KeyUp:
		mode.editor.Motion(File.Up)
	case KeyDown:
		mode.editor.Motion(File.Down)
	case KeyHome:
		mode.editor.Motion(File.LineStart)
	case KeyEnd:
		mode.editor.Motion(File.LineEnd)
	case KeyPageUp:
		mode.editor.Motion(Paragraph(Backward))
	case KeyPageDown:
		mode.editor.Motion(Paragraph(Forward))
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
		mode.editor.Motion(File.Left)
	case 'f':
		mode.editor.Motion(File.Right)
	case 'a':
		mode.editor.Motion(File.Up)
	case 's':
		mode.editor.Motion(File.Down)
	case 'D':
		mode.editor.Motion(File.LineStart)
	case 'F':
		mode.editor.Motion(File.LineEnd)
	case 'A':
		mode.editor.Motion(File.FileStart)
	case 'S':
		mode.editor.Motion(File.FileEnd)
	case 'e':
		mode.editor.Motion(Word(Backward))
	case 'r':
		mode.editor.Motion(Word(Forward))
	case 'q':
		mode.editor.Motion(Paragraph(Backward))
	case 'w':
		mode.editor.Motion(Paragraph(Forward))
	case 'g':
		mode.editor.Copy()
		mode.editor.SwitchMode(mode.editor.Command)
	case 'j':
		mode.editor.Copy()
		mode.editor.Change(File.Delete)
		mode.editor.SwitchMode(mode.editor.Command)
	}
	if err != nil {
		return fmt.Errorf("error handling rune %v: %w", rune, err)
	}
	return nil
}

// Render renders select mode.
func (mode *Select) Render(content *Content) error {
	mode.editor.File.Render(content, true)
	content.Color = ColorYellow
	content.Prompt = ""
	return nil
}
