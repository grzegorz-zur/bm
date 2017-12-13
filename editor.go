package bm

import (
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Editor struct {
	*File
	Mode
	Normal  *Normal
	Input   *Input
	restart chan struct{}
	exit    bool
}

func (editor *Editor) ApplyFileOp(op FileOp) {
	*editor.File = op(*editor.File)
}

func Open(path string) (editor *Editor, err error) {
	file, err := Read(path)
	editor = &Editor{
		File:    &file,
		restart: make(chan struct{}),
	}
	editor.Normal = &Normal{
		Editor: editor,
	}
	editor.Input = &Input{
		Editor: editor,
	}
	editor.SwitchMode(editor.Normal)
	return
}

func (editor *Editor) SwitchMode(mode Mode) {
	editor.Mode = mode
}

func (editor *Editor) Quit() (err error) {
	editor.exit = true
	return
}

func (editor *Editor) Close() (err error) {
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
	editor.signals()

	for !editor.exit {
		err = editor.display()
		if err != nil {
			err = errors.Wrap(err, "display failed")
			editor.Quit()
		}
		err = editor.listen()
		if err != nil {
			err = errors.Wrap(err, "event poll failed")
			editor.Quit()
		}
	}
	return
}

func (editor *Editor) signals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGCONT)
	go func() {
		for s := range signals {
			var err error
			switch s {
			case syscall.SIGCONT:
				err = editor.cont()
			}
			if err != nil {
				log.Fatalf("signal handling failure", err)
			}
		}
	}()
}

func (editor *Editor) listen() (err error) {
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

func (editor *Editor) display() (err error) {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	width, height := tb.Size()
	size := Size{Lines: height, Cols: width}
	bounds := Bounds{Right: size.Cols - 1, Bottom: size.Lines - 1}
	cursor, err := editor.Mode.Display(bounds)
	tb.SetCursor(cursor.Col, cursor.Line)
	tb.Flush()
	return
}

func (editor *Editor) Stop() (err error) {
	tb.Close()
	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	if err != nil {
		err = errors.Wrap(err, "editor stop failed")
		return
	}
	p.Signal(syscall.SIGSTOP)
	<-editor.restart
	return
}

func (editor *Editor) cont() (err error) {
	err = tb.Init()
	if err != nil {
		err = errors.Wrap(err, "editor continue failed")
		return
	}
	editor.restart <- struct{}{}
	return
}
