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

func OpenFile(path string) (file *File, err error) {
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

func (file *File) Resize(size Size) {
	file.Window.Bottom = file.Window.Top + size.Lines
	file.Window.Right = file.Window.Left + size.Cols
	if file.Position.Line > file.Window.Bottom {
		file.Position.Line = file.Window.Bottom
	}
	if file.Position.Col > file.Window.Right {
		file.Position.Col = file.Window.Right
	}
	return
}

func (file *File) Display(position Position) (cursor Position) {
	for line := file.Window.Top; line <= file.Window.Bottom; line++ {
		if line >= len(file.Data) {
			break
		}
		runes := file.Data[line]
		absLine := position.Line + line
		for col := file.Window.Left; col <= file.Window.Right; col++ {
			if col >= len(runes) {
				break
			}
			symbol := runes[col]
			absCol := position.Col + col
			tb.SetCell(absCol, absLine, symbol, tb.ColorDefault, tb.ColorDefault)
		}
	}
	cursor = file.Position
	return
}

func (file *File) InsertAfter(r rune) {
	p := &file.Position
	if len(file.Data) <= p.Line {
		line := []rune{}
		file.Data = append(file.Data, line)
	}
	if len(file.Data[p.Line]) <= p.Col {
		file.Data[p.Line] = append(file.Data[p.Line], r)
	} else {
		line := file.Data[p.Line]
		line = append(line[:p.Col], append([]rune{r}, line[p.Col:]...)...)
		file.Data[p.Line] = line
	}
	p.Col += 1
}
