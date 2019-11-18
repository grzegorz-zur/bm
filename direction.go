package main

// Direction	indicates forward or backward.
type Direction int

const (
	// Forward direction.
	Forward Direction = 1
	// Backword direction.
	Backward Direction = -1
)

// Value returns increment for calculation
func (d Direction) Value() int {
	return int(d)
}
