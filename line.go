package main

// Line represents line of characters.
type Line []rune

// DeleteRune deletes a character at the a column.
func (l Line) DeleteRune(c int) Line {
	if c >= len(l) {
		return l
	}
	ln := make(Line, 0, len(l)-1)
	ln = append(ln, l[:c]...)
	ln = append(ln, l[c+1:]...)
	return ln
}

// DeletePreviousRune deletes a character left from the a column.
func (l Line) DeletePreviousRune(c int) Line {
	if c == 0 || c > len(l) {
		return l
	}
	ln := make(Line, 0, len(l)-1)
	ln = append(ln, l[:c-1]...)
	ln = append(ln, l[c:]...)
	return ln
}

// Expand expands line to a number of columns.
//
// The line is filled with space characters.
func (l Line) Expand(c int) Line {
	if c < len(l) {
		return l
	}
	ln := make(Line, 0, c)
	ln = append(ln, l...)
	for i := len(ln); i < cap(ln); i++ {
		ln = append(ln, ' ')
	}
	return ln
}

// InsertRune inserts a character at a column.
func (l Line) InsertRune(c int, r rune) Line {
	le := l.Expand(c)
	ln := make(Line, 0, len(l)+1)
	ln = append(ln, le[:c]...)
	ln = append(ln, r)
	ln = append(ln, le[c:]...)
	return ln
}

// AppendRune appends rune at the end.
func (l Line) AppendRune(r rune) Line {
	return l.InsertRune(len(l), r)
}

// Split splits line at a column.
func (l Line) Split(c int) (Line, Line) {
	le := l.Expand(c)
	l1 := make(Line, 0, c)
	l2 := make(Line, 0, len(le)-c)
	l1 = append(l1, le[:c]...)
	l2 = append(l2, le[c:]...)
	return l1, l2
}

// String makes string of a line.
func (l Line) String() string {
	return string(l)
}
