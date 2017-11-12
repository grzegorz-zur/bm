package bm

func ExpandLines(ls [][]rune, l int) (lines [][]rune) {
	if l < len(ls) {
		return ls
	}
	lines = make([][]rune, l)
	for i := range ls {
		lines[i] = ls[i]
	}
	return
}

func ExpandLine(l []rune, c int) (line []rune) {
	if c < len(l) {
		return l
	}
	line = make([]rune, c)
	for i := range l {
		line[i] = l[i]
	}
	for i := len(l); i < len(line); i++ {
		line[i] = ' '
	}
	return
}
