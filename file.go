package bm

import (
	"bufio"
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

var (
	ErrNoFile             = errors.New("no file")
	eol                   = []byte("\n")
	perm      os.FileMode = 0644
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
	if file == nil {
		return
	}
	file.Position = move(*file)
}

func (file *File) Change(op Change) {
	if file == nil {
		return
	}
	*(file) = op(*file)
	file.Archive()
}

func (file *File) Archive() {
	if file == nil {
		return
	}
	file.History.Archive(file.Lines, file.Position)
}

func (file *File) SwitchVersion(dir Direction) {
	if file == nil {
		return
	}
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

func (file *File) ReloadIfModified() (modified bool, err error) {
	if file == nil {
		return modified, ErrNoFile
	}
	modified, err = file.Modified()
	if err != nil {
		return
	}
	if !modified {
		return
	}
	err = file.Reload()
	return
}

func (file *File) Modified() (modified bool, err error) {
	if file == nil {
		return modified, ErrNoFile
	}
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
	if file == nil {
		return ErrNoFile
	}
	err = file.Load()
	if err != nil {
		err = errors.Wrapf(err, "error reloading file: %s", file.Path)
	}
	file.Archive()
	return
}

func (file *File) Load() (err error) {
	if file == nil {
		return ErrNoFile
	}
	flags := os.O_RDWR | os.O_CREATE
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
	if file == nil {
		return ErrNoFile
	}
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
	for _, runes := range file.Lines {
		line := string(runes)
		bytes := []byte(line)
		f.Write(bytes)
		f.Write(eol)
	}
	return
}

func (file *File) Render(display *Display, bounds Bounds) (cursor Position, err error) {
	if file == nil {
		return cursor, ErrNoFile
	}
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
	if file == nil {
		return ErrNoFile
	}
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
	if file == nil {
		return
	}
	w := &file.Window
	w.Bottom = w.Top + size.Lines
	w.Right = w.Left + size.Cols
	return
}

func (file *File) scroll() {
	if file == nil {
		return
	}
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

func (file File) DeleteRune() File {
	file.Lines = file.Lines.DeleteRune(file.Position)
	return file
}

func (file File) DeletePreviousRune() File {
	p := file.Position
	if p.Col == 0 {
		return file
	}
	file.Lines = file.Lines.DeletePreviousRune(p)
	file.Position = Position{Line: p.Line, Col: p.Col - 1}
	return file
}

func InsertRune(r rune) Change {
	return func(file File) File {
		p := file.Position
		file.Lines = file.Lines.InsertRune(p, r)
		file.Position = Position{Line: p.Line, Col: p.Col + 1}
		return file
	}
}

func (file File) DeleteLine() File {
	file.Lines = file.Lines.DeleteLine(file.Position.Line)
	return file
}

func (file File) Split() File {
	file.Lines = file.Lines.Split(file.Position)
	file.Position = Position{Line: file.Position.Line + 1}
	return file
}
