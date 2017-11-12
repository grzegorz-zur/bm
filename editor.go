package bm

import (
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"log"
)

type Editor struct {
	Size
	*File
	Mode
	Normal *Normal
	Input  *Input
	exit   bool
}

type Position struct {
	Line int
	Col  int
}

type Size struct {
	Lines int
	Cols  int
}

type Bounds struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

func (editor *Editor) ApplyFileOp(op FileOp) {
	*editor.File = op(*editor.File)
}

func Open(path string) (editor *Editor, err error) {
	file, err := Read(path)
	editor = &Editor{
		File: file,
	}
	editor.Normal = &Normal{
		Editor: editor,
	}
	editor.Input = &Input{
		Editor: editor,
	}
	editor.Switch(editor.Normal)
	return
}

func (editor *Editor) Run() (err error) {
	err = tb.Init()
	if err != nil {
		err = errors.Wrap(err, "editor init failed")
		return
	}
	defer tb.Close()

	defer editor.Close()

	width, height := tb.Size()
	size := Size{Cols: width, Lines: height}
	editor.Resize(size)

	for !editor.exit {
		editor.Display()
		err = editor.Listen()
		if err != nil {
			err = errors.Wrap(err, "event poll failed")
			editor.Quit()
		}
		editor.Scroll()
	}
	return
}

func (editor *Editor) Switch(mode Mode) (err error) {
	editor.Mode = mode
	return
}

func (editor *Editor) Listen() (err error) {
	event := tb.PollEvent()
	switch event.Type {
	case tb.EventKey:
		err = editor.Key(event)
	case tb.EventError:
		err = errors.Wrap(event.Err, "event poll failed")
	default:
		log.Printf("%+v\n", event)
	}
	return
}

func (editor *Editor) Resize(size Size) {
	editor.Size = size
	editor.File.Resize(size)
	return
}

func (editor *Editor) Display() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	cursor := editor.File.Display(Position{0, 0})
	tb.SetCursor(cursor.Col, cursor.Line)
	tb.Flush()
	return
}

func (editor *Editor) Quit() (err error) {
	editor.exit = true
	return
}

func (editor *Editor) Close() (err error) {
	return
}
