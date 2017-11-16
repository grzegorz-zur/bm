package bm

func DeleteRune(f File) (file File) {
	file = f
	file.Data = DeleteRuneLines(f.Data, f.Position)
	return
}

func DeleteRuneLines(ls [][]rune, p Position) (lines [][]rune) {
	if p.Line >= len(ls) {
		return ls
	}
	l := ls[p.Line]
	if p.Col >= len(l) {
		return ls
	}
	line := DeleteRuneLine(l, p.Col)
	lines = make([][]rune, len(ls))
	for i := range ls {
		lines[i] = ls[i]
	}
	lines[p.Line] = line
	return
}

func DeleteRuneLine(l []rune, c int) (line []rune) {
	if c >= len(l) {
		return l
	}
	line = append(line, l[:c]...)
	line = append(line, l[c+1:]...)
	return
}
