package bm

type Line []rune

func (l Line) DeleteRune(c int) (line Line) {
	if c >= len(l) {
		return l
	}
	line = append(line, l[:c]...)
	line = append(line, l[c+1:]...)
	return
}

func (l Line) Expand(c int) (line Line) {
	if c < len(l) {
		return l
	}
	line = make(Line, c)
	for i := range l {
		line[i] = l[i]
	}
	for i := len(l); i < len(line); i++ {
		line[i] = ' '
	}
	return
}

func (l Line) InsertRune(c int, r rune) (line Line) {
	el := l.Expand(c)
	line = append(line, el[:c]...)
	line = append(line, r)
	line = append(line, el[c:]...)
	return
}

func (l Line) Split(c int) (l1, l2 Line) {
	el := l.Expand(c)
	l1 = append(l1, el[:c]...)
	l2 = append(l2, el[c:]...)
	return
}
