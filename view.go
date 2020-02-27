package main

// View holds data to be displayed on screen.
type View struct {
	Size      Size
	Visible   bool
	Select    bool
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
	content := make([][]rune, size.Lines)
	selection := make([][]bool, size.Lines)
	for line := 0; line < size.Lines; line++ {
		content[line] = make([]rune, size.Columns)
		selection[line] = make([]bool, size.Columns)
	}
	return &View{
		Size:      size,
		Content:   content,
		Selection: selection,
	}
}

// Clear clears view.
func (view *View) Clear() {
	for line := 0; line < view.Size.Lines; line++ {
		for column := 0; column < view.Size.Columns; column++ {
			view.Content[line][column] = 0
			view.Selection[line][column] = false
		}
	}
	view.Visible = false
	view.Select = false
	view.Position = Position{}
	view.Color = ColorNone
	view.Status = ""
	view.Prompt = ""
	view.Cursor = CursorNone
}
