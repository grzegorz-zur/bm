package main

// Area describes a rectangular area.
//
// All values are inclusive.
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
		L: a.B - a.T + 1,
		C: a.R - a.L + 1,
	}
}

// SplitHorizontal splits area horizontally.
//
// Positive argument defines size of the top area, negative defines size of the bottom area.
func (a Area) SplitHorizontal(l int) (Area, Area) {
	s := 0
	if l >= 0 {
		s = a.T + l
	} else {
		s = a.B + l
	}
	ah, al := a, a
	ah.B = s
	al.T = s + 1
	return ah, al
}
