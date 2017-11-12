package bm

type Move func(f File) (p Position)

func MoveOp(m Move) FileOp {
	return func(f File) (file File) {
		file = f
		file.Position = m(f)
		return
	}
}

func Left(f File) (p Position) {
	p = f.Position
	if p.Col > 0 {
		p.Col--
	}
	return
}

func Right(f File) (p Position) {
	p = f.Position
	p.Col++
	return
}

func Up(f File) (p Position) {
	p = f.Position
	if p.Line > 0 {
		p.Line--
	}
	return
}

func Down(f File) (p Position) {
	p = f.Position
	p.Line++
	return
}
