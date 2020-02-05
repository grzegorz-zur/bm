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
		mode.editor.Motion(File.Left)
	case KeyRight:
		mode.editor.Motion(File.Right)
	case KeyUp:
		mode.editor.Motion(File.Up)
	case KeyDown:
		mode.editor.Motion(File.Down)
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
func (mode *Command) Rune(rune rune) error {
	var err error
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
		mode.editor.Motion(File.Left)
	case 'f':
		mode.editor.Motion(File.Right)
	case 'a':
		mode.editor.Motion(File.Up)
	case 's':
		mode.editor.Motion(File.Down)
	case 'e':
		mode.editor.Motion(Word(Backward))
	case 'r':
		mode.editor.Motion(Word(Forward))
	case 'q':
		mode.editor.Motion(Paragraph(Backward))
	case 'w':
		mode.editor.Motion(Paragraph(Forward))
	case 'z':
		mode.editor.SwitchVersion(Backward)
	case 'x':
		mode.editor.SwitchVersion(Forward)
	case 'c':
		mode.editor.SwitchFile(Backward)
	case 'v':
		mode.editor.SwitchFile(Forward)
	case 'j':
		mode.editor.Change(File.DeleteRune)
	case 'J':
		mode.editor.Change(File.DeleteLine)
	case 'g':
		mode.editor.SwitchMode(mode.editor.Select)
	case 'h':
		mode.editor.PasteBlock()
	case 'H':
		mode.editor.PasteInline()
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
func (mode *Command) Render(content *Content) error {
	if !mode.editor.Empty() {
		mode.editor.File.Render(content, false)
	}
	content.Color = ColorGreen
	content.Prompt = ""
	return nil
}
