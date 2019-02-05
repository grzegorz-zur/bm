package bm

type Lines []Line

func (ls Lines) DeleteRune(pos Position) (lines Lines) {
	if pos.Line >= len(ls) {
		return ls
	}
	l := ls[pos.Line]
	line := l.DeleteRune(pos.Col)
	lines = make(Lines, len(ls))
	for i := range ls {
		lines[i] = ls[i]
	}
	lines[pos.Line] = line
	return
}

func (ls Lines) DeletePreviousRune(pos Position) (lines Lines) {
	if pos.Line >= len(ls) {
		return ls
	}
	l := ls[pos.Line]
	line := l.DeletePreviousRune(pos.Col)
	lines = make(Lines, len(ls))
	for i := range ls {
		lines[i] = ls[i]
	}
	lines[pos.Line] = line
	return
}

func (ls Lines) DeleteLine(l int) (lines Lines) {
	if l >= len(ls) {
		return ls
	}
	lines = make(Lines, len(ls)-1)
	for i := 0; i < l; i++ {
		lines[i] = ls[i]
	}
	for j := l + 1; j < len(ls); j++ {
		lines[j-1] = ls[j]
	}
	return
}

func (ls Lines) Expand(l int) (lines Lines) {
	if l < len(ls) {
		return ls
	}
	lines = make(Lines, l)
	for i := range ls {
		lines[i] = ls[i]
	}
	return
}

func (ls Lines) InsertRune(pos Position, r rune) (lines Lines) {
	els := ls.Expand(pos.Line + 1)
	l := els[pos.Line]
	line := l.InsertRune(pos.Col, r)
	lines = make(Lines, len(els))
	for i := range els {
		lines[i] = els[i]
	}
	lines[pos.Line] = line
	return
}

func (ls Lines) Split(pos Position) (lines Lines) {
	els := ls.Expand(pos.Line + 1)
	l := els[pos.Line]
	line1, line2 := l.Split(pos.Col)
	lines = make(Lines, len(els)+1)
	for i := 0; i < pos.Line; i++ {
		lines[i] = els[i]
	}
	lines[pos.Line] = line1
	lines[pos.Line+1] = line2
	for i := pos.Line + 1; i < len(ls); i++ {
		lines[i+1] = els[i]
	}
	return
}
