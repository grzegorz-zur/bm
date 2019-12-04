package main

// Input is a mode for typing.
type Input struct {
	*Editor
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
	case KeyEscape:
		m.SwitchMode(m.Editor.Command)
	case KeyLeft:
		m.Motion(File.Left)
	case KeyRight:
		m.Motion(File.Right)
	case KeyUp:
		m.Motion(File.Up)
	case KeyDown:
		m.Motion(File.Down)
	case KeyPageUp:
		m.Motion(Paragraph(Backward))
	case KeyPageDown:
		m.Motion(Paragraph(Forward))
	case KeyTab:
		m.Change(InsertRune('\t'))
	case KeyEnter:
		m.Change(File.Split)
	case KeyBackspace:
		m.Change(File.DeletePreviousRune)
	case KeyDelete:
		m.Change(File.DeleteRune)
	}
	return nil
}

// Rune handles rune input.
func (m *Input) Rune(r rune) error {
	m.Change(InsertRune(r))
	return nil
}

func (m *Input) Render(cnt *Content) error {
	m.File.Render(cnt)
	cnt.Color = ColorRed
	return nil
}
