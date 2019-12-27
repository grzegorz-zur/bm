package main

// Position indicates line and column.
//
// The values start at zero.
type Position struct {
	// Line.
	L int
	// Column.
	C int
}

// Less checks if first position is strictly smaller than the second position.
func Less(a, b Position) bool {
	switch {
	case a.L < b.L:
		return true
	case a.L > b.L:
		return false
	default:
		return a.C < b.C
	}
}

// Sort sorts the positions.
func Sort(a, b Position) (Position, Position) {
	if Less(a, b) {
		return a, b
	}
	return b, a
}

// Between checks if position is between positions.
func Between(p, a, b Position) bool {
	a, b = Sort(a, b)
	return (p == a || Less(a, p)) && (p == b || Less(p, b))
}

// Min returns the smaller position.
func Min(a, b Position) Position {
	p, _ := Sort(a, b)
	return p
}
