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
	Base    string
	keys    chan tb.Event
	pause   chan struct{}
	unpause chan struct{}
	quit    chan struct{}
	done    chan struct{}
}

func New(display *Display, path string, files []string) (editor *Editor) {
	editor = &Editor{
		Display: display,
		Base:    path,
		keys:    make(chan tb.Event),
		pause:   make(chan struct{}, 1),
		unpause: make(chan struct{}),
		quit:    make(chan struct{}, 1),
		done:    make(chan struct{}),
	}
	editor.Modes.Command = &Command{
		Editor: editor,
	}
	editor.Modes.Input = &Input{
		Editor: editor,
	}
	editor.Modes.Switch = &Switch{
		Editor: editor,
	}
	for _, file := range files {
		editor.Open(file)
	}
	if editor.File == nil {
		editor.SwitchMode(editor.Switch)
	} else {
		editor.Next(Forward)
		editor.SwitchMode(editor.Command)
	}
	return
}

func (editor *Editor) Start() {
	go editor.signals()
	go editor.listen()
	go editor.run()
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

func (editor *Editor) Pause() {
	editor.pause <- struct{}{}
	return
}

func (editor *Editor) Quit() {
	editor.quit <- struct{}{}
	return
}

func (editor *Editor) Wait() {
	<-editor.done
}

func (editor *Editor) signals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGCONT, syscall.SIGTERM)
	for signal := range signals {
		switch signal {
		case syscall.SIGCONT:
			editor.unpause <- struct{}{}
		case syscall.SIGTERM:
			editor.quit <- struct{}{}
		}
	}
}

func (editor *Editor) listen() {
	for {
		editor.keys <- tb.PollEvent()
	}
}

func (editor *Editor) run() {
	defer close(editor.done)
	err := editor.Display.Init()
	if err != nil {
		err = errors.Wrap(err, "display init failed")
		return
	}
	defer editor.Display.Close()

	for {
		if editor.File == nil {
			editor.SwitchMode(editor.Switch)
		}
		err = editor.render()
		if err != nil {
			err = errors.Wrap(err, "render failed")
			log.Println(err)
			return
		}
		select {
		case event := <-editor.keys:
			err = editor.Key(event)
			if err != nil {
				err = errors.Wrap(err, "key handling failed")
				log.Println(err)
				return
			}
		case <-editor.pause:
			err = editor.background()
			if err != nil {
				err = errors.Wrap(err, "pause handling failed")
				log.Println(err)
				return
			}
		case <-editor.quit:
			return
		}
	}
}

func (editor *Editor) background() (err error) {
	editor.Display.Close()
	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	if err != nil {
		err = errors.Wrap(err, "background failed")
		return
	}
	p.Signal(syscall.SIGSTOP)
	<-editor.unpause
	err = editor.Display.Init()
	if err != nil {
		err = errors.Wrap(err, "display init failed")
		return
	}
	return
}

func (editor *Editor) render() (err error) {
	editor.Display.Clear(tb.ColorDefault, tb.ColorDefault)
	width, height := editor.Display.Size()
	size := Size{Lines: height, Cols: width}
	bounds := Bounds{Right: size.Cols - 1, Bottom: size.Lines - 1}
	cursor, err := editor.Mode.Render(editor.Display, bounds)
	editor.Display.SetCursor(cursor.Col, cursor.Line)
	editor.Display.Flush()
	return
}
