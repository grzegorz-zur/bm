package main

// Lines represents lines of text.
type Lines []Line

// Slice slices lines from position to position.
func (ls Lines) Slice(a, b Position) Lines {
	s := make(Lines, 0, b.L-a.L+1)
	for l := a.L; l <= b.L; l++ {
		if l < len(ls) {
			s = append(s, ls[l])
		} else {
			s = append(s, Line{})
		}
	}
	if len(s) > 0 {
		s[0] = s[0].Suffix(a.C)
		s[len(s)-1] = s[len(s)-1].Prefix(b.C + 1)
	}
	return s
}

// Delete deletes content between positions.
func (ls Lines) Delete(a, b Position) Lines {
	if len(ls) == 0 {
		return ls
	}
	lsl := len(ls) - 1
	ma := min(a.L, lsl)
	mb := min(b.L, lsl)
	lsnl := lsl - (ma - mb)
	lsn := make(Lines, 0, lsnl)
	lsn = append(lsn, ls[:ma]...)
	ln := Line{}
	if a.L < len(ls) {
		p := ls[a.L].Prefix(a.C)
		ln = append(ln, p...)
	}
	if b.L < len(ls) {
		s := ls[b.L].Suffix(b.C + 1)
		ln = append(ln, s...)
	}
	if a.L < len(ls) {
		lsn = append(lsn, ln)
	}
	lsn = append(lsn, ls[mb+1:]...)
	return lsn
}

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

// InsertBlock insert lines at a position as a block.
func (ls Lines) InsertBlock(p Position, lsi Lines) Lines {
	lse := ls.Expand(p.L + 1)
	lsn := make(Lines, 0, len(lse)+len(lsi))
	lsn = append(lsn, lse[:p.L]...)
	lsn = append(lsn, lsi...)
	lsn = append(lsn, lse[p.L:]...)
	return lsn
}

// InsertInline inserts lines at a position inline.
func (ls Lines) InsertInline(p Position, lsi Lines) Lines {
	if len(lsi) == 0 {
		return ls
	}
	lse := ls.Expand(p.L + 1)
	lsn := make(Lines, 0, len(lse)+len(lsi)-1)
	lsn = append(lsn, lse[:p.L]...)

	lfe := lse[p.L].Expand(p.C)
	lfi := lsi[0]
	lfn := make(Line, 0, len(lfe)+len(lfi))
	lfn = append(lfn, lfe[:p.C]...)
	lfn = append(lfn, lfi...)
	lsn = append(lsn, lfn)

	if len(lsi) > 2 {
		lsn = append(lsn, lsi[1:len(lsi)-1]...)
	}

	lle := lfe[p.C:]
	lli := lsi[len(lsi)-1]
	lln := make(Line, 0, len(lle)+len(lli))
	lln = append(lln, lli...)
	lln = append(lln, lle...)
	lsn = append(lsn, lln)

	lsn = append(lsn, lse[p.L+1:]...)
	return lsn
}

// AtText checks if position is at text.
func (ls Lines) AtText(p Position) bool {
	if p.L >= len(ls) {
		return false
	}
	if p.C >= len(ls[p.L]) {
		return false
	}
	return true
}
