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
func NewContent(size Size) *Content {
	runes := make([][]rune, size.L)
	marks := make([][]bool, size.L)
	for l := 0; l < size.L; l++ {
		runes[l] = make([]rune, size.C)
		marks[l] = make([]bool, size.C)
	}
	return &Content{
		Size:  size,
		Runes: runes,
		Marks: marks,
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
	cnt.Position = Position{}
	cnt.Color = ColorNone
	cnt.Status = ""
	cnt.Prompt = ""
	cnt.Cursor = CursorNone
}
