package main

type Visual interface {
	Render(*Display, Bounds) (Position, error)
}
