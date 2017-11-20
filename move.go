package bm

type Move func(f File) (p Position)

func (f File) Left() (p Position) {
	p = f.Position
	if p.Col > 0 {
		p.Col--
	}
	return
}

func (f File) Right() (p Position) {
	p = f.Position
	p.Col++
	return
}

func (f File) Up() (p Position) {
	p = f.Position
	if p.Line > 0 {
		p.Line--
	}
	return
}

func (f File) Down() (p Position) {
	p = f.Position
	p.Line++
	return
}
