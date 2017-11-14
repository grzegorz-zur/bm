package bm

func InsertRune(r rune) FileOp {
	return func(f File) (file File) {
		file = f
		p := f.Position
		file.Data = InsertRuneLines(f.Data, p, r)
		file.Position = Position{Line: p.Line, Col: p.Col + 1}
		return
	}
}

func InsertRuneLines(ls [][]rune, p Position, r rune) (lines [][]rune) {
	els := ExpandLines(ls, p.Line+1)
	l := els[p.Line]
	line := InsertRuneLine(l, p.Col, r)
	lines = make([][]rune, len(els))
	for i := range els {
		lines[i] = els[i]
	}
	lines[p.Line] = line
	return
}

func InsertRuneLine(l []rune, c int, r rune) (line []rune) {
	el := ExpandLine(l, c)
	line = append(line, el[:c]...)
	line = append(line, r)
	line = append(line, el[c:]...)
	return
}
