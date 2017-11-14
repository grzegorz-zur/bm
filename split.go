package bm

func Split(f File) (file File) {
	file = f
	p := f.Position
	file.Data = SplitLines(f.Data, p)
	file.Position = Position{Line: p.Line + 1}
	return
}

func SplitLines(ls [][]rune, p Position) (lines [][]rune) {
	els := ExpandLines(ls, p.Line+1)
	l := els[p.Line]
	line1, line2 := SplitLine(l, p.Col)
	lines = make([][]rune, len(els)+1)
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

func SplitLine(l []rune, c int) (line1, line2 []rune) {
	el := ExpandLine(l, c)
	line1 = append(line1, el[:c]...)
	line2 = append(line2, el[c:]...)
	return
}
