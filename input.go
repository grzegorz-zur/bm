package main

// Input is a mode for typing.
type Input struct {
	editor *Editor
}

// Show updates mode when switched to.
func (m *Input) Show() error {
	return nil
}

// Hide updates mode when switched from.
func (m *Input) Hide() error {
	return nil
}

// Key handles input events.
func (m *Input) Key(k Key) error {
	switch k {
	case KeyLeft:
		m.editor.Motion(File.Left)
	case KeyRight:
		m.editor.Motion(File.Right)
	case KeyUp:
		m.editor.Motion(File.Up)
	case KeyDown:
		m.editor.Motion(File.Down)
	case KeyPageUp:
		m.editor.Motion(Paragraph(Backward))
	case KeyPageDown:
		m.editor.Motion(Paragraph(Forward))
	case KeyTab:
		m.editor.Change(InsertRune('\t'))
	case KeyEnter:
		m.editor.Change(File.Split)
	case KeyBackspace:
		m.editor.Change(File.DeletePreviousRune)
	case KeyDelete:
		m.editor.Change(File.DeleteRune)
	case KeyCtrlSpace:
		m.editor.SwitchMode(m.editor.Command)
	}
	return nil
}

// Rune handles rune input.
func (m *Input) Rune(r rune) error {
	m.editor.Change(InsertRune(r))
	return nil
}

// Render renders mode to the screen.
func (m *Input) Render(cnt *Content) error {
	m.editor.File.Render(cnt, false)
	cnt.Color = ColorRed
	return nil
}
