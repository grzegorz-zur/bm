package bm

func Split(f File) (file File) {
	file = f
	file.Data = SplitLines(f.Data, f.Position)
	file.Position = Position{Line: f.Position.Line + 1}
	return
}

func SplitLines(ls [][]rune, p Position) (lines [][]rune) {
	ls = ExpandLines(ls, p.Line+1)
	l := ls[p.Line]
	l1, l2 := SplitLine(l, p.Col)
	lines = make([][]rune, len(ls)+1)
	for i := 0; i < p.Line; i++ {
		lines[i] = ls[i]
	}
	lines[p.Line] = l1
	lines[p.Line+1] = l2
	for i := p.Line + 1; i < len(ls); i++ {
		lines[i+1] = ls[i]
	}
	return
}

func SplitLine(l []rune, c int) (l1, l2 []rune) {
	l = ExpandLine(l, c)
	l1 = append(l1, l[:c]...)
	l2 = append(l2, l[c:]...)
	return
}
