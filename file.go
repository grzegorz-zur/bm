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

// Position calculates position.
func (file *File) Position(location int) (position Position, ok bool) {
	l, c := 0, 0
	for loc, rune := range file.content {
		if loc == location {
			return Position{l, c}, true
		}
		c++
		if rune == EOL {
			l++
			c = 0
		}
	}
	return Position{l, c}, false
}

// Location calculates location.
func (file *File) Location(line, col int) (location int, ok bool) {
	l, c := 0, 0
	for loc, rune := range file.content {
		location = loc
		if l == line && c == col {
			return location, true
		} else if l == line && rune == EOL {
			return location, false
		}
		c++
		if rune == EOL {
			l++
			c = 0
		}
	}
	return len(file.content), false
}

// Archive makes a record in history.
func (file *File) Archive() {
	if file == nil {
		return
	}
	file.history.Archive(file.content, file.location)
}

// SwitchVersion switches between versions from history.
func (file *File) SwitchVersion(dir Direction) {
	if file == nil {
		return
	}
	file.content, file.location = file.history.Switch(dir)
}

// Select sets selection position.
func (file *File) Select() {
	if file == nil {
		return
	}
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
	_, size := file.rune(to)
	to += size
	return file.content[from:to], from, to
}

// Render renders file content.
func (file *File) Render(view *View, selection bool) error {
	position, _ := file.Position(file.location)
	file.area = file.area.Resize(view.Size).Shift(position)
	for line := file.area.T; line < file.area.B; line++ {
		rline := line - file.area.T
		for col := file.area.L; col < file.area.R; col++ {
			rcol := col - file.area.L
			location, ok := file.Location(line, col)
			if ok {
				rune, _ := utf8.DecodeRuneInString(file.content[location:])
				view.Content[rline][rcol] = rune
				view.Selection[rline][rcol] = selection && file.selected(location)
			} else {
				break
			}
		}
	}
	view.Position = Position{
		L: position.L - file.area.T,
		C: position.C - file.area.L,
	}
	view.Status = fmt.Sprintf("%s %d:%d", file.Path, position.L+1, position.C+1)
	view.Cursor = CursorContent
	return nil
}

func (file *File) selected(location int) bool {
	return file.mark <= location && location < file.location ||
		file.location < location && location <= file.mark
}

func (file *File) last() (rune rune, size int) {
	return utf8.DecodeLastRuneInString(file.content[:file.location])
}

func (file *File) current() (rune rune, size int) {
	return utf8.DecodeRuneInString(file.content[file.location:])
}

func (file *File) rune(location int) (rune rune, size int) {
	return utf8.DecodeRuneInString(file.content[location:])
}
