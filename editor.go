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
	Mode
	Normal *Normal
	Input  *Input
	*File
	Files
	restart chan struct{}
	exit    bool
}

func (editor *Editor) ApplyFileOp(op FileOp) {
	*editor.File = op(*editor.File)
}

func New() (editor *Editor) {
	editor = &Editor{
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

func (editor *Editor) New() {
	file := NewFile()
	editor.File = &file
	editor.Files = editor.Files.Add(&file)
}

func (editor *Editor) Open(path string) (err error) {
	file, err := Read(path)
	if err != nil {
		return
	}
	editor.File = &file
	editor.Files = editor.Files.Add(&file)
	return
}

func (editor *Editor) Next(d Direction) {
	editor.SwitchFile(editor.Files.Next(editor.File, d))
}

func (editor *Editor) Close() {
	file := editor.File
	editor.File = editor.Files.Next(file, Forward)
	editor.Files = editor.Files.Remove(file)
	if editor.Files.Empty() {
		editor.File = nil
	}
}

func (editor *Editor) SwitchFile(f *File) {
	editor.File = f
}

func (editor *Editor) SwitchMode(mode Mode) {
	editor.Mode = mode
}

func (editor *Editor) Quit() (err error) {
	editor.exit = true
	return
}

func (editor *Editor) Run() (err error) {
	err = tb.Init()
	if err != nil {
		err = errors.Wrap(err, "editor init failed")
		return
	}
	defer tb.Close()

	editor.signals()

	for !editor.exit {
		editor.init()
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

func (editor *Editor) init() {
	if editor.File == nil {
		editor.New()
	}
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
