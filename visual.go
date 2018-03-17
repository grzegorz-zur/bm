package bm

type Visual interface {
	Render(*Display, Bounds) (Position, error)
}
