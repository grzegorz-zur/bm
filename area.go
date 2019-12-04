package main

// Area describes a rectangular area.
type Area struct {
	// Top.
	T int
	// Bottom.
	B int
	// Left.
	L int
	// Right.
	R int
}

// Size calculates the number of lines and columns in the area.
func (a Area) Size() Size {
	return Size{
		L: a.B - a.T,
		C: a.R - a.L,
	}
}

// Resize resizes the area.
func (a Area) Resize(s Size) Area {
	return Area{
		T: a.T,
		B: a.T + s.L,
		L: a.L,
		R: a.L + s.C,
	}
}

// Shift shifts area to include position.
func (a Area) Shift(p Position) Area {
	s := a.Size()
	switch {
	case p.L < a.T:
		a.T = p.L
	case p.L >= a.B:
		a.T += p.L - a.B + 1
	}
	switch {
	case p.C < a.L:
		a.L = p.C
	case p.C >= a.R:
		a.L += p.C - a.R + 1
	}
	return a.Resize(s)
}
