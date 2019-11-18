package main

// Visual is an object that renders in terminal.
type Visual interface {
	Render(*Display, Area) (Position, error)
}
