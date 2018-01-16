package bm

type Direction int

const (
	Forward  Direction = 1
	Backward Direction = -1
)

func (dir Direction) Value() int {
	return int(dir)
}
