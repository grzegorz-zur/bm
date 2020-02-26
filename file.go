package main

import (
	"fmt"
	"time"
	"unicode/utf8"
)

// File represents open file.
type File struct {
	Path     string
	content  string
	location int
	mark     int
	area     Area
	time     time.Time
	changed  bool
	history  *History
}

const (
	// Tab is tabulation rune.
	Tab = '\t'
	// EOL is end of line rune.
	EOL = '\n'
)

// Archive makes a record in history.
func (file *File) Archive() {
	file.history.Archive(file.content, file.location)
}

// SwitchVersion switches between versions from history.
func (file *File) SwitchVersion(dir Direction) {
	file.content, file.location = file.history.Switch(dir)
}

// Select sets selection position.
func (file *File) Select() {
	file.mark = file.location
}

// Copy returns selected content.
func (file *File) Copy() string {
	content, _, _ := file.copy()
	return content
}

// Cut cuts selected content.
func (file *File) Cut() string {
	content, from, to := file.copy()
	file.Remove(from, to)
	return content
}

func (file *File) copy() (content string, from, to int) {
	from, to = file.mark, file.location
	if from > to {
		from, to = to, from
	}
	_, size := file.next(to)
	to += size
	return file.content[from:to], from, to
}

// Render renders file content.
func (file *File) Render(view *View, selection bool) error {
	position := file.position()
	file.area = file.area.Resize(view.Size).Shift(position)
	line, column := 0, 0
	for location, rune := range file.content {
		if file.area.Contains(Position{line, column}) {
			rline := line - file.area.Top
			rcolumn := column - file.area.Left
			view.Content[rline][rcolumn] = rune
			view.Selection[rline][rcolumn] = selection && file.selected(location)
		}
		column++
		if rune == EOL {
			line++
			column = 0
		}
	}
	view.Position = Position{
		Line:   position.Line - file.area.Top,
		Column: position.Column - file.area.Left,
	}
	view.Status = fmt.Sprintf("%s %d:%d", file.Path, position.Line+1, position.Column+1)
	view.Cursor = CursorContent
	return nil
}

func (file *File) position() Position {
	line, column := 0, 0
	for location, rune := range file.content {
		if location == file.location {
			break
		}
		column++
		if rune == EOL {
			line++
			column = 0
		}
	}
	return Position{line, column}
}

func (file *File) selected(location int) bool {
	return file.mark <= location && location < file.location ||
		file.location < location && location <= file.mark
}

func (file *File) last() (rune rune, size int) {
	return file.previous(file.location)
}

func (file *File) current() (rune rune, size int) {
	return file.next(file.location)
}

func (file *File) previous(location int) (rune rune, size int) {
	return utf8.DecodeLastRuneInString(file.content[:location])
}

func (file *File) next(location int) (rune rune, size int) {
	return utf8.DecodeRuneInString(file.content[location:])
}
