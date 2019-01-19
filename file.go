package bm

import (
	"bufio"
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

type File struct {
	Path string
	Time time.Time
	Lines
	Position
	Window Bounds
	*History
}

type Change func(File) File

func (file *File) Move(move Move) {
	file.Position = move(*file)
}

func (file *File) Change(op Change) {
	*(file) = op(*file)
	file.Archive()
}

func (file *File) Archive() {
	file.History.Archive(file.Lines, file.Position)
}

func (file *File) SwitchVersion(dir Direction) {
	file.Lines, file.Position = file.History.Switch(dir)
}

func Open(path string) (file File, err error) {
	file = File{
		Path:    path,
		History: &History{},
	}
	err = file.Load()
	if err != nil {
		return
	}
	file.Archive()
	return
}

func (file *File) Modified() (modified bool, err error) {
	stat, err := os.Stat(file.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		err = errors.Wrapf(err, "error checking file: %s", file.Path)
	}
	modified = stat.ModTime() != file.Time
	return
}

func (file *File) Reload() (err error) {
	err = file.Load()
	if err != nil {
		err = errors.Wrapf(err, "error reloading file: %s", file.Path)
	}
	file.Archive()
	return
}

func (file *File) Load() (err error) {
	flags := os.O_RDWR | os.O_CREATE
	perm := os.ModePerm
	f, err := os.OpenFile(file.Path, flags, perm)
	if err != nil {
		err = errors.Wrapf(err, "error opening file: %s", file.Path)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			err = errors.Wrapf(err, "error closing file: %s", file.Path)
			log.Print(err)
		}
	}()
	err = file.update()
	if err != nil {
		err = errors.Wrapf(err, "error updating file information: %s", file.Path)
		log.Print(err)
	}
	s := bufio.NewScanner(f)
	file.Lines = nil
	for s.Scan() {
		err = s.Err()
		if err != nil {
			err = errors.Wrapf(err, "error reading file: %s", file.Path)
			return
		}
		line := s.Text()
		runes := []rune(line)
		file.Lines = append(file.Lines, runes)
	}
	return
}

func (file *File) Write() (err error) {
	f, err := os.Create(file.Path)
	defer func() {
		err := f.Close()
		if err != nil {
			err = errors.Wrapf(err, "error closing file: %s", file.Path)
			log.Print(err)
		}
	}()
	defer func() {
		err := file.update()
		if err != nil {
			err = errors.Wrapf(err, "error updating file information: %s", file.Path)
			log.Print(err)
		}
	}()
	if err != nil {
		err = errors.Wrapf(err, "error writing file: %s", file.Path)
		return
	}
	for i, runes := range file.Lines {
		line := string(runes)
		if i+1 < len(file.Lines) {
			line += "\n"
		}
		bytes := []byte(line)
		f.Write(bytes)
	}
	return
}

func (file *File) Render(display *Display, bounds Bounds) (cursor Position, err error) {
	file.scroll()
	size := bounds.Size()
	file.size(size)
	w := file.Window
	for line := w.Top; line <= w.Bottom; line++ {
		if line >= len(file.Lines) {
			break
		}
		runes := file.Lines[line]
		screenLine := bounds.Top + line - w.Top
		for col := w.Left; col <= w.Right; col++ {
			if col >= len(runes) {
				break
			}
			symbol := runes[col]
			screenCol := bounds.Left + col - w.Left
			display.SetCell(screenCol, screenLine, symbol, tb.ColorDefault, tb.ColorDefault)
		}
	}
	p := file.Position
	cursor.Line = p.Line - w.Top
	cursor.Col = p.Col - w.Left
	return
}

func (file *File) update() (err error) {
	stat, err := os.Stat(file.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		err = errors.Wrapf(err, "error checking file: %s", file.Path)
	}
	file.Time = stat.ModTime()
	return
}

func (file *File) size(size Size) {
	w := &file.Window
	w.Bottom = w.Top + size.Lines
	w.Right = w.Left + size.Cols
	return
}

func (file *File) scroll() {
	p := file.Position
	w := &file.Window
	height := w.Bottom - w.Top
	width := w.Right - w.Left

	switch {
	case p.Line < w.Top:
		w.Top = p.Line
		w.Bottom = w.Top + height
	case p.Line > w.Bottom:
		w.Bottom = p.Line
		w.Top = w.Bottom - height
	}

	switch {
	case p.Col < w.Left:
		w.Left = p.Col
		w.Right = w.Left + width
	case p.Col > w.Right:
		w.Right = p.Col
		w.Left = w.Right - width
	}
}

func (f File) DeleteRune() File {
	f.Lines = f.Lines.DeleteRune(f.Position)
	return f
}

func (f File) DeletePreviousRune() File {
	p := f.Position
	if p.Col == 0 {
		return f
	}
	f.Lines = f.Lines.DeletePreviousRune(p)
	f.Position = Position{Line: p.Line, Col: p.Col - 1}
	return f
}

func InsertRune(r rune) Change {
	return func(f File) File {
		p := f.Position
		f.Lines = f.Lines.InsertRune(p, r)
		f.Position = Position{Line: p.Line, Col: p.Col + 1}
		return f
	}
}

func (f File) DeleteLine() File {
	f.Lines = f.Lines.DeleteLine(f.Position.Line)
	return f
}

func (f File) Split() File {
	f.Lines = f.Lines.Split(f.Position)
	f.Position = Position{Line: f.Position.Line + 1}
	return f
}
