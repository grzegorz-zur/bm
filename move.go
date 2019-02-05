package bm

type Move func(file File) (pos Position)

func (file File) Left() (pos Position) {
	pos = file.Position
	if pos.Col > 0 {
		pos.Col--
	}
	return
}

func (file File) Right() (pos Position) {
	pos = file.Position
	pos.Col++
	return
}

func (file File) Up() (pos Position) {
	pos = file.Position
	if pos.Line > 0 {
		pos.Line--
	}
	return
}

func (file File) Down() (pos Position) {
	pos = file.Position
	pos.Line++
	return
}
