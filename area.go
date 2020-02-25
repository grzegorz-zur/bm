package main

// Area describes a rectangular area.
type Area struct {
	// Top.
	Top int
	// Bottom.
	Bottom int
	// Left.
	Left int
	// Right.
	Right int
}

// Size calculates the number of lines and columns in the area.
func (area Area) Size() Size {
	return Size{
		Lines:   area.Bottom - area.Top,
		Columns: area.Right - area.Left,
	}
}

// Resize resizes the area.
func (area Area) Resize(size Size) Area {
	return Area{
		Top:    area.Top,
		Bottom: area.Top + size.Lines,
		Left:   area.Left,
		Right:  area.Left + size.Columns,
	}
}

// Shift shifts area to include position.
func (area Area) Shift(position Position) Area {
	size := area.Size()
	switch {
	case position.Line < area.Top:
		area.Top = position.Line
	case position.Line >= area.Bottom:
		area.Top += position.Line - area.Bottom + 1
	}
	switch {
	case position.Column < area.Left:
		area.Left = position.Column
	case position.Column >= area.Right:
		area.Left += position.Column - area.Right + 1
	}
	return area.Resize(size)
}
