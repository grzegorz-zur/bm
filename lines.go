package main

type Lines []Line

func (ls Lines) DeleteRune(pos Position) (lines Lines) {
	if pos.Line >= len(ls) {
		return ls
	}
	l := ls[pos.Line]
	line := l.DeleteRune(pos.Col)
	lines = make(Lines, 0, len(ls))
	lines = append(lines, ls...)
	lines[pos.Line] = line
	return
}

func (ls Lines) DeletePreviousRune(pos Position) (lines Lines) {
	if pos.Line >= len(ls) {
		return ls
	}
	l := ls[pos.Line]
	line := l.DeletePreviousRune(pos.Col)
	lines = make(Lines, 0, len(ls))
	lines = append(lines, ls...)
	lines[pos.Line] = line
	return
}

func (ls Lines) DeleteLine(l int) (lines Lines) {
	if l >= len(ls) {
		return ls
	}
	lines = make(Lines, 0, len(ls)-1)
	lines = append(lines, ls[:l]...)
	lines = append(lines, ls[l+1:]...)
	return
}

func (ls Lines) Expand(l int) (lines Lines) {
	if l < len(ls) {
		return ls
	}
	lines = make(Lines, 0, l)
	lines = append(lines, ls...)
	for i := len(lines); i < cap(lines); i++ {
		lines = append(lines, Line{})
	}
	return
}

func (ls Lines) InsertRune(pos Position, r rune) (lines Lines) {
	els := ls.Expand(pos.Line + 1)
	l := els[pos.Line]
	line := l.InsertRune(pos.Col, r)
	lines = make(Lines, 0, len(els))
	lines = append(lines, els[:pos.Line]...)
	lines = append(lines, line)
	lines = append(lines, els[pos.Line+1:]...)
	return
}

func (ls Lines) Split(pos Position) (lines Lines) {
	els := ls.Expand(pos.Line + 1)
	l := els[pos.Line]
	line1, line2 := l.Split(pos.Col)
	lines = make(Lines, 0, len(els)+1)
	lines = append(lines, els[:pos.Line]...)
	lines = append(lines, line1, line2)
	lines = append(lines, els[pos.Line+1:]...)
	return
}
