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
	*Display
	Modes
	Files
	Base string
	wait chan struct{}
	exit bool
}

func New(display *Display, path string) (editor *Editor) {
	editor = &Editor{
		Display: display,
		Base:    path,
		wait:    make(chan struct{}),
	}
	editor.Modes.Normal = &Normal{
		Editor: editor,
	}
	editor.Modes.Input = &Input{
		Editor: editor,
	}
	editor.Modes.Switch = &Switch{
		Editor: editor,
	}
	editor.SwitchMode(editor.Normal)
	return
}

func (editor *Editor) Open(path string) (err error) {
	return editor.Files.Open(editor.Base, path)
}

func (editor *Editor) Write() (err error) {
	return editor.Files.Write(editor.Base)
}

func (editor *Editor) WriteAll() (err error) {
	return editor.Files.WriteAll(editor.Base)
}

func (editor *Editor) Quit() (err error) {
	editor.exit = true
	return
}

func (editor *Editor) Run() (err error) {
	err = editor.Display.Init()
	if err != nil {
		err = errors.Wrap(err, "editor init failed")
		return
	}
	defer editor.Display.Close()

	editor.signals()

	for !editor.exit {
		err = editor.display()
		if err != nil {
			err = errors.Wrap(err, "display failed")
			log.Println(err)
			editor.Quit()
		}
		err = editor.listen()
		if err != nil {
			err = errors.Wrap(err, "event poll failed")
			log.Println(err)
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
				log.Fatalf("signal handling failure %v", err)
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
	editor.Display.Clear(tb.ColorDefault, tb.ColorDefault)
	width, height := editor.Display.Size()
	size := Size{Lines: height, Cols: width}
	bounds := Bounds{Right: size.Cols - 1, Bottom: size.Lines - 1}
	cursor, err := editor.Mode.Render(editor.Display, bounds)
	editor.Display.SetCursor(cursor.Col, cursor.Line)
	editor.Display.Flush()
	return
}

func (editor *Editor) Stop() (err error) {
	editor.Display.Close()
	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	if err != nil {
		err = errors.Wrap(err, "editor stop failed")
		return
	}
	p.Signal(syscall.SIGSTOP)
	<-editor.wait
	return
}

func (editor *Editor) cont() (err error) {
	err = editor.Display.Init()
	if err != nil {
		err = errors.Wrap(err, "editor continue failed")
		return
	}
	editor.wait <- struct{}{}
	return
}
