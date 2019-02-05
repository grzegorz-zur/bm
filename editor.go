package bm

import (
	tb "github.com/nsf/termbox-go"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	TickInterval = 200 * time.Millisecond
)

type Editor struct {
	*Display
	Modes
	Files
	keys    chan tb.Event
	check   chan struct{}
	pause   chan struct{}
	unpause chan struct{}
	quit    chan struct{}
	done    chan struct{}
}

func New(display *Display, files []string) (editor *Editor) {
	editor = &Editor{
		Display: display,
		keys:    make(chan tb.Event),
		check:   make(chan struct{}),
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
		editor.SwitchFile(Forward)
		editor.SwitchMode(editor.Command)
	}
	return
}

func (editor *Editor) Start() {
	go editor.signals()
	go editor.listen()
	go editor.tick()
	go editor.run()
}

func (editor *Editor) Check() {
	editor.check <- struct{}{}
}

func (editor *Editor) Pause() {
	editor.pause <- struct{}{}
}

func (editor *Editor) Quit() {
	editor.quit <- struct{}{}
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

func (editor *Editor) tick() {
	for {
		time.Sleep(TickInterval)
		editor.check <- struct{}{}
	}
}

func (editor *Editor) run() {
	defer close(editor.done)
	err := editor.Display.Init()
	if err != nil {
		report(err)
		return
	}
	defer editor.Display.Close()

	for {
		if editor.File == nil {
			editor.SwitchMode(editor.Switch)
		}
		err = editor.render()
		report(err)
		select {
		case event := <-editor.keys:
			err = editor.Key(event)
			report(err)
		case <-editor.check:
			if !editor.Empty() {
				_, err = editor.ReloadIfModified()
				report(err)
			}
		case <-editor.pause:
			err = editor.background()
			report(err)
		case <-editor.quit:
			return
		}
	}
}

func report(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (editor *Editor) background() (err error) {
	editor.Display.Close()
	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err != nil {
		err = errors.Wrapf(err, "error pausing process %+v", process)
		return
	}
	process.Signal(syscall.SIGSTOP)
	<-editor.unpause
	err = editor.Display.Init()
	if err != nil {
		err = errors.Wrap(err, "error initializing display")
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
	if err != nil {
		err = errors.Wrap(err, "error rendering editor")
	}
	editor.Display.SetCursor(cursor.Col, cursor.Line)
	editor.Display.Flush()
	return
}
