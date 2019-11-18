package main

// Lines represents lines of text.
type Lines []Line

// DeleteRune deletes a character at a position.
func (ls Lines) DeleteRune(p Position) Lines {
	if p.L >= len(ls) {
		return ls
	}
	l := ls[p.L]
	ln := l.DeleteRune(p.C)
	lsn := make(Lines, 0, len(ls))
	lsn = append(lsn, ls...)
	lsn[p.L] = ln
	return lsn
}

// DeletePreviousRune deletes a character left to a position.
func (ls Lines) DeletePreviousRune(p Position) Lines {
	if p.L >= len(ls) {
		return ls
	}
	l := ls[p.L]
	ln := l.DeletePreviousRune(p.C)
	lsn := make(Lines, 0, len(ls))
	lsn = append(lsn, ls...)
	lsn[p.L] = ln
	return lsn
}

// DeleteLine deletes a line at a row.
func (ls Lines) DeleteLine(l int) Lines {
	if l >= len(ls) {
		return ls
	}
	lsn := make(Lines, 0, len(ls)-1)
	lsn = append(lsn, ls[:l]...)
	lsn = append(lsn, ls[l+1:]...)
	return lsn
}

// Expand expands lines to a number of rows.
func (ls Lines) Expand(l int) Lines {
	if l < len(ls) {
		return ls
	}
	lsn := make(Lines, 0, l)
	lsn = append(lsn, ls...)
	for i := len(lsn); i < cap(lsn); i++ {
		lsn = append(lsn, Line{})
	}
	return lsn
}

// InsertRune inserts a character at a position.
func (ls Lines) InsertRune(p Position, r rune) Lines {
	lse := ls.Expand(p.L + 1)
	l := lse[p.L]
	ln := l.InsertRune(p.C, r)
	lsn := make(Lines, 0, len(lse))
	lsn = append(lsn, lse[:p.L]...)
	lsn = append(lsn, ln)
	lsn = append(lsn, lse[p.L+1:]...)
	return lsn
}

// Split splits lines at a position.
func (ls Lines) Split(p Position) Lines {
	lse := ls.Expand(p.L + 1)
	l := lse[p.L]
	ln1, ln2 := l.Split(p.C)
	lsn := make(Lines, 0, len(lse)+1)
	lsn = append(lsn, lse[:p.L]...)
	lsn = append(lsn, ln1, ln2)
	lsn = append(lsn, lse[p.L+1:]...)
	return lsn
}
