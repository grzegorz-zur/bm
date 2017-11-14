package bm

import (
	"bufio"
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"log"
	"os"
)

type File struct {
	Path     string
	Data     [][]rune
	Window   Bounds
	Position Position
}

type FileOp func(f File) (file File)

func Read(path string) (file *File, err error) {
	file = &File{
		Path: path,
	}
	f, err := os.Open(path)
	if err != nil {
		err = errors.Wrapf(err, "cannot open file: %s", path)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			err = errors.Wrapf(err, "cannot close file: %s", path)
			log.Print(err)
		}
	}()
	s := bufio.NewScanner(f)
	for s.Scan() {
		err = s.Err()
		if err != nil {
			err = errors.Wrapf(err, "cannot read file: %s", path)
			return
		}
		line := s.Text()
		runes := []rune(line)
		file.Data = append(file.Data, runes)
	}
	return
}

func (file *File) Write() (err error) {
	f, err := os.Create(file.Path)
	if err != nil {
		err = errors.Wrapf(err, "cannot write file: %s", file.Path)
		return
	}
	for i, runes := range file.Data {
		line := string(runes)
		if i+1 < len(file.Data) {
			line += "\n"
		}
		bytes := []byte(line)
		f.Write(bytes)
	}
	return
}

func (file *File) Scroll() {
	p := &file.Position
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

func (file *File) Resize(size Size) {
	p := &file.Position
	w := &file.Window
	w.Bottom = w.Top + size.Lines - 1
	w.Right = w.Left + size.Cols - 1
	if p.Line > w.Bottom {
		p.Line = w.Bottom
	}
	if p.Col > w.Right {
		p.Col = w.Right
	}
	return
}

func (file *File) Display(position Position) (cursor Position) {
	w := file.Window
	for line := w.Top; line <= w.Bottom; line++ {
		if line >= len(file.Data) {
			break
		}
		runes := file.Data[line]
		screenLine := position.Line + line - w.Top
		for col := w.Left; col <= w.Right; col++ {
			if col >= len(runes) {
				break
			}
			symbol := runes[col]
			screenCol := position.Col + col - w.Left
			tb.SetCell(screenCol, screenLine, symbol, tb.ColorDefault, tb.ColorDefault)
		}
	}
	p := file.Position
	cursor.Line = p.Line - w.Top
	cursor.Col = p.Col - w.Left
	return
}

func (file *File) Delete() {
	switch {
	case file.empty():
		return
	case file.emptyLine():
		file.DeleteLine()
	case file.emptyChar():
		return
	default:
		file.DeleteChar()
	}
}

func (file *File) DeleteChar() {
	p := &file.Position
	line := &file.Data[p.Line]
	rest := (*line)[p.Col+1:]
	*line = append((*line)[:p.Col], rest...)
}

func (file *File) DeleteLine() {
	p := &file.Position
	data := &file.Data
	rest := (*data)[p.Line+1:]
	*data = append(*data, rest...)
}

func (file *File) empty() bool {
	p := &file.Position
	return p.Line >= len(file.Data) ||
		p.Col >= len(file.Data[p.Line])
}

func (file *File) emptyLine() bool {
	p := &file.Position
	line := &file.Data[p.Line]
	return len(*line) == 0
}

func (file *File) emptyChar() bool {
	p := &file.Position
	return p.Col >= len(file.Data[p.Line])
}

func (file *File) extend() {
	file.extendLine()
	file.extendCol()
}

func (file *File) extendLine() {
	p := &file.Position
	data := &file.Data
	for p.Line >= len(*data) {
		*data = append(*data, []rune{})
	}
}

func (file *File) extendCol() {
	p := &file.Position
	line := &file.Data[p.Line]
	for p.Col >= len(*line) {
		*line = append(*line, ' ')
	}
}
