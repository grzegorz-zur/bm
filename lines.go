package bm

type Lines []Line

func (ls Lines) DeleteRune(p Position) (lines Lines) {
	if p.Line >= len(ls) {
		return ls
	}
	l := ls[p.Line]
	if p.Col >= len(l) {
		return ls
	}
	line := l.DeleteRune(p.Col)
	lines = make(Lines, len(ls))
	for i := range ls {
		lines[i] = ls[i]
	}
	lines[p.Line] = line
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

func (ls Lines) InsertRune(p Position, r rune) (lines Lines) {
	els := ls.Expand(p.Line + 1)
	l := els[p.Line]
	line := l.InsertRune(p.Col, r)
	lines = make(Lines, len(els))
	for i := range els {
		lines[i] = els[i]
	}
	lines[p.Line] = line
	return
}

func (ls Lines) Split(p Position) (lines Lines) {
	els := ls.Expand(p.Line + 1)
	l := els[p.Line]
	line1, line2 := l.Split(p.Col)
	lines = make(Lines, len(els)+1)
	for i := 0; i < p.Line; i++ {
		lines[i] = els[i]
	}
	lines[p.Line] = line1
	lines[p.Line+1] = line2
	for i := p.Line + 1; i < len(ls); i++ {
		lines[i+1] = els[i]
	}
	return
}
