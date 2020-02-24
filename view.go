package main

// View holds data to be displayed on screen.
type View struct {
	Size      Size
	Content   [][]rune
	Selection [][]bool
	Position  Position
	Color     Color
	Status    string
	Prompt    string
	Cursor    Cursor
}

// Cursor represents type of cursor.
type Cursor int

// Cursor types.
const (
	CursorNone Cursor = iota
	CursorContent
	CursorPrompt
)

// NewView creates a new content of a given size.
func NewView(size Size) *View {
	content := make([][]rune, size.L)
	selection := make([][]bool, size.L)
	for line := 0; line < size.L; line++ {
		content[line] = make([]rune, size.C)
		selection[line] = make([]bool, size.C)
	}
	return &View{
		Size:      size,
		Content:   content,
		Selection: selection,
	}
}

// Clear clears view.
func (view *View) Clear() {
	for line := 0; line < view.Size.L; line++ {
		for col := 0; col < view.Size.C; col++ {
			view.Content[line][col] = 0
			view.Selection[line][col] = false
		}
	}
	view.Position = Position{}
	view.Color = ColorNone
	view.Status = ""
	view.Prompt = ""
	view.Cursor = CursorNone
}
