package main

// Content holds data to be displayed on screen.
type Content struct {
	Size     Size
	Runes    [][]rune
	Marks    [][]bool
	Position Position
	Color    Color
	Status   string
	Prompt   string
	Cursor   Cursor
}

// Cursor represents type of cursor.
type Cursor int

// Cursor types.
const (
	CursorNone Cursor = iota
	CursorContent
	CursorPrompt
)

// NewContent creates a new content of a given size.
func NewContent(s Size) *Content {
	rs := make([][]rune, s.L)
	ms := make([][]bool, s.L)
	for l := 0; l < s.L; l++ {
		rs[l] = make([]rune, s.C)
		ms[l] = make([]bool, s.C)
	}
	return &Content{
		Size:  s,
		Runes: rs,
		Marks: ms,
	}
}

// Clear clears content.
func (cnt *Content) Clear() {
	for l := 0; l < cnt.Size.L; l++ {
		for c := 0; c < cnt.Size.C; c++ {
			cnt.Runes[l][c] = 0
			cnt.Marks[l][c] = false
		}
	}
}
