package bm

type Visual interface {
	Display(Bounds) (Position, error)
}
