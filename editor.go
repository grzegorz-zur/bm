package bm

import (
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
)

type Editor struct {
	File *File
}

func Init() (editor *Editor, err error) {
	editor = &Editor{
		File: &File{},
	}
	return
}

func (editor *Editor) Run() (err error) {
	err = tb.Init()
	if err != nil {
		err = errors.Wrap(err, "editor init failed")
		return
	}
	defer tb.Close()

	for {
		err = editor.Listen()
		if err != nil {
			err = errors.Wrap(err, "event poll failed")
			return
		}
		err = editor.Display()
		if err != nil {
			err = errors.Wrap(err, "display failed")
			return
		}
	}
}

func (editor *Editor) Listen() (err error) {
	event := tb.PollEvent()
	switch event.Type {
	case tb.EventKey:
		err = editor.Key(event)
	case tb.EventError:
		err = errors.Wrap(event.Err, "event poll failed")
	}
	return
}

func (editor *Editor) Key(event tb.Event) (err error) {
	editor.File.Key(event)
	return
}

func (editor *Editor) Display() (err error) {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	columns, lines := tb.Size()
	column, line, err := editor.File.Display(0, 0, columns, lines)
	if err != nil {
		err = errors.Wrapf(err, "display of file %s failed", editor.File.Path)
	}
	tb.SetCursor(column, line)
	tb.Flush()
	return
}
