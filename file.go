package main

import (
	"errors"
	"fmt"
	tb "github.com/nsf/termbox-go"
	"os"
	"time"
)

var (
	// ErrNoFile indicates missing file.
	ErrNoFile             = errors.New("no file")
	eol                   = []byte("\n")
	perm      os.FileMode = 0644
)

// File represents open file.
type File struct {
	Path string
	Time time.Time
	Lines
	Position
	Window Area
	*History
}

// Motion applies motion to a file.
func (f *File) Motion(m Motion) {
	if f == nil {
		return
	}
	f.Position = m(*f)
}

// Change applies change to a file.
func (f *File) Change(c Change) {
	if f == nil {
		return
	}
	*(f) = c(*f)
	f.Archive()
}

// Archive makes a record in history.
func (f *File) Archive() {
	if f == nil {
		return
	}
	f.History.Archive(f.Lines, f.Position)
}

// SwitchVersion switches between versions from history.
func (f *File) SwitchVersion(dir Direction) {
	if f == nil {
		return
	}
	f.Lines, f.Position = f.History.Switch(dir)
}

// Render renders the file contents to display.
func (f *File) Render(d *Display, a Area) (Position, error) {
	if f == nil {
		return Position{}, ErrNoFile
	}
	f.scroll()
	s := a.Size()
	f.size(s)
	w := f.Window
	for l := w.T; l <= w.B; l++ {
		if l >= len(f.Lines) {
			break
		}
		runes := f.Lines[l]
		sl := a.T + l - w.T
		for c := w.L; c <= w.R; c++ {
			if c >= len(runes) {
				break
			}
			symbol := runes[c]
			sc := a.L + c - w.L
			d.SetCell(sc, sl, symbol, tb.ColorDefault, tb.ColorDefault)
		}
	}
	p := f.Position
	crs := Position{
		L: p.L - w.T,
		C: p.C - w.L,
	}
	return crs, nil
}

func (f *File) update() error {
	if f == nil {
		return ErrNoFile
	}
	stat, err := os.Stat(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error checking file %s: %w", f.Path, err)
	}
	f.Time = stat.ModTime()
	return nil
}

func (f *File) size(s Size) {
	if f == nil {
		return
	}
	w := &f.Window
	w.B = w.T + s.L - 1
	w.R = w.L + s.C - 1
	return
}

func (f *File) scroll() {
	if f == nil {
		return
	}
	p := f.Position
	w := &f.Window
	s := w.Size()

	switch {
	case p.L < w.T:
		w.T = p.L
		w.B = w.T + s.L - 1
	case p.L > w.B:
		w.B = p.L
		w.T = w.B - s.L + 1
	}

	switch {
	case p.C < w.L:
		w.L = p.C
		w.R = w.L + s.C - 1
	case p.C > w.R:
		w.R = p.C
		w.L = w.R - s.C + 1
	}
}
