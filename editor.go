package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// TickInterval designates time between file modification checks.
	TickInterval = 200 * time.Millisecond
)

// Editor represents the editor.
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

// New creates new editor.
func New(d *Display, fs []string) *Editor {
	e := &Editor{
		Display: d,
		keys:    make(chan tb.Event),
		check:   make(chan struct{}),
		pause:   make(chan struct{}, 1),
		unpause: make(chan struct{}),
		quit:    make(chan struct{}, 1),
		done:    make(chan struct{}),
	}
	e.Modes.Command = &Command{
		Editor: e,
	}
	e.Modes.Input = &Input{
		Editor: e,
	}
	e.Modes.Switch = &Switch{
		Editor: e,
	}
	for _, f := range fs {
		e.Open(f)
	}
	if e.Empty() {
		e.SwitchMode(e.Switch)
	} else {
		e.SwitchFile(Forward)
		e.SwitchMode(e.Command)
	}
	return e
}

// Start starts the editor.
func (e *Editor) Start() {
	go e.signals()
	go e.listen()
	go e.tick()
	go e.run()
}

// SendKey sends event to the editor.
func (e *Editor) SendKey(ev tb.Event) {
	e.keys <- ev
}

// Check signals file modification check.
func (e *Editor) Check() {
	e.check <- struct{}{}
}

// Pause pauses the editor.
func (e *Editor) Pause() {
	e.pause <- struct{}{}
}

// Quit quits the editor.
func (e *Editor) Quit() {
	e.quit <- struct{}{}
}

// Wait waits for the editor to finish.
func (e *Editor) Wait() {
	<-e.done
}

func (e *Editor) signals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGCONT, syscall.SIGTERM)
	for signal := range signals {
		switch signal {
		case syscall.SIGCONT:
			e.unpause <- struct{}{}
		case syscall.SIGTERM:
			e.quit <- struct{}{}
		}
	}
}

func (e *Editor) listen() {
	for {
		e.keys <- tb.PollEvent()
	}
}

func (e *Editor) tick() {
	for {
		time.Sleep(TickInterval)
		e.check <- struct{}{}
	}
}

func (e *Editor) run() {
	defer close(e.done)
	err := e.Display.Init()
	if err != nil {
		report(err)
		return
	}
	defer e.Display.Close()

	for {
		if e.Empty() {
			e.SwitchMode(e.Switch)
		}
		err = e.render()
		report(err)
		select {
		case event := <-e.keys:
			err = e.Key(event)
			report(err)
		case <-e.check:
			if !e.Empty() {
				_, err = e.ReloadIfModified()
				report(err)
			}
		case <-e.pause:
			err = e.background()
			report(err)
		case <-e.quit:
			return
		}
	}
}

func report(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (e *Editor) background() error {
	e.Display.Close()
	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("error pausing process %+v: %w", process, err)
	}
	process.Signal(syscall.SIGSTOP)
	<-e.unpause
	err = e.Display.Init()
	if err != nil {
		return fmt.Errorf("error initializing display: %w", err)
	}
	return nil
}

func (e *Editor) render() error {
	e.Display.Clear(tb.ColorDefault, tb.ColorDefault)
	width, height := e.Display.Size()
	size := Size{L: height, C: width}
	area := Area{R: size.C - 1, B: size.L - 1}
	cursor, err := e.Mode.Render(e.Display, area)
	if err != nil {
		return fmt.Errorf("error rendering editor: %w", err)
	}
	e.Display.SetCursor(cursor.C, cursor.L)
	e.Display.Flush()
	return nil
}
