package main

import (
	"fmt"
	"os"
	"time"
)

var (
	eol              = []byte("\n")
	perm os.FileMode = 0644
)

// File represents open file.
type File struct {
	Path string
	Lines
	*History
	pos     Position
	sel     Position
	area    Area
	time    time.Time
	changed bool
}

// Motion applies motion to a file.
func (f *File) Motion(m Motion) {
	if f == nil {
		return
	}
	f.pos = m(*f)
}

// Change applies change to a file.
func (f *File) Change(c Change) {
	if f == nil {
		return
	}
	*(f) = c(*f)
	f.changed = true
	f.Archive()
}

// Select sets selection position.
func (f *File) Select() {
	if f == nil {
		return
	}
	f.sel = f.pos
}

// Archive makes a record in history.
func (f *File) Archive() {
	if f == nil {
		return
	}
	f.History.Archive(f.Lines, f.pos)
}

// SwitchVersion switches between versions from history.
func (f *File) SwitchVersion(dir Direction) {
	if f == nil {
		return
	}
	f.Lines, f.pos = f.History.Switch(dir)
}

// Selection returns selected lines.
func (f *File) Selection() Lines {
	s, e := Sort(f.sel, f.pos)
	return f.Slice(s, e)
}

// Render renders file content.
func (f *File) Render(cnt *Content, mark bool) error {
	f.area = f.area.Resize(cnt.Size).Shift(f.pos)
	for l := f.area.T; l < f.area.B; l++ {
		rl := l - f.area.T
		for c := f.area.L; c < f.area.R; c++ {
			rc := c - f.area.L
			if l < len(f.Lines) {
				line := f.Lines[l]
				if c < len(line) {
					cnt.Runes[rl][rc] = line[c]
					if mark {
						p := Position{l, c}
						cnt.Marks[rl][rc] = f.marked(p)
					}
				}
			}
		}
	}
	cnt.Position = Position{
		L: f.pos.L - f.area.T,
		C: f.pos.C - f.area.L,
	}
	cnt.Status = fmt.Sprintf("%s %d:%d", f.Path, f.pos.L+1, f.pos.C+1)
	cnt.Cursor = CursorContent
	return nil
}

func (f *File) marked(p Position) bool {
	return Between(p, f.sel, f.pos)
}
