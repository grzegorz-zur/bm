package main

import (
	"fmt"
	"github.com/gdamore/tcell"
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
	Modes
	Files
	newScreen NewScreen
	screen    tcell.Screen
	view      *View
	content   string
	events    chan tcell.Event
	check     chan struct{}
	pause     chan struct{}
	unpause   chan struct{}
	quit      chan struct{}
	done      chan struct{}
}

// NewScreen creates a new screen.
type NewScreen func() (tcell.Screen, error)

// New creates new editor.
func New(ns NewScreen, fs []string) *Editor {
	e := &Editor{
		newScreen: ns,
		view:      NewView(Size{}),
		events:    make(chan tcell.Event),
		check:     make(chan struct{}),
		pause:     make(chan struct{}, 1),
		unpause:   make(chan struct{}),
		quit:      make(chan struct{}, 1),
		done:      make(chan struct{}),
	}
	e.Modes.Command = &Command{
		editor: e,
	}
	e.Modes.Input = &Input{
		editor: e,
	}
	e.Modes.Select = &Select{
		editor: e,
	}
	e.Modes.Switch = &Switch{
		editor: e,
	}
	for _, f := range fs {
		e.Open(f)
	}
	e.SwitchFile(Forward)
	e.SwitchMode(e.Command)
	return e
}

// Start starts the editor.
func (e *Editor) Start() error {
	var err error
	e.screen, err = e.newScreen()
	if err != nil {
		return fmt.Errorf("error creating screen: %w", err)
	}
	err = e.screen.Init()
	if err != nil {
		return fmt.Errorf("error initializing screen: %w", err)
	}
	go e.signals()
	go e.listen()
	go e.tick()
	go e.run()
	return nil
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

// Copy copies selection to buffer.
func (e *Editor) Copy() {
	e.content = e.File.Copy()
	e.SwitchMode(e.Command)
}

// Cut cuts selection to buffer.
func (e *Editor) Cut() {
	e.content = e.File.Cut()
	e.SwitchMode(e.Command)
}

// Paste pastes buffer.
func (e *Editor) Paste() {
	e.Insert(e.content)
}

// LineAbove starts line above current line.
func (e *Editor) LineAbove() {
	e.MoveLineStart()
	e.Insert(string(EOL))
	e.MoveUp()
	e.SwitchMode(e.Input)
}

// LineBelow starts line below current line.
func (e *Editor) LineBelow() {
	e.MoveLineEnd()
	e.Insert(string(EOL))
	e.SwitchMode(e.Input)
}

func (e *Editor) listen() {
	for {
		e.events <- e.screen.PollEvent()
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
	render := true
	for {
		if render {
			err := e.render()
			e.report(err)
		}
		select {
		case ev := <-e.events:
			err := e.handle(ev)
			e.report(err)
			render = true
		case <-e.check:
			if e.Empty() {
				render = false
			} else {
				r, err := e.Read(false)
				render = r
				e.report(err)
			}
		case <-e.pause:
			err := e.background()
			e.report(err)
		case <-e.quit:
			e.close()
			return
		}
	}
}

func (e *Editor) handle(event tcell.Event) error {
	switch tevent := event.(type) {
	case *tcell.EventKey:
		if tevent.Key() == tcell.KeyRune {
			return e.Mode.Rune(tevent.Rune())
		}
		key, ok := keymap[tevent.Key()]
		if ok {
			return e.Mode.Key(key)
		}
	case *tcell.EventResize:
		width, height := tevent.Size()
		if height > 0 {
			height--
		}
		size := Size{height, width}
		e.view = NewView(size)
	}
	return nil
}

func (e *Editor) report(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (e *Editor) background() error {
	e.screen.Fini()
	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("error pausing process %+v: %w", process, err)
	}
	process.Signal(syscall.SIGSTOP)
	<-e.unpause
	e.screen, err = e.newScreen()
	if err != nil {
		return fmt.Errorf("error creating screen: %w", err)
	}
	err = e.screen.Init()
	if err != nil {
		return fmt.Errorf("error initializing screen: %w", err)
	}
	return nil
}

func (e *Editor) close() {
	e.screen.Fini()
}

func (e *Editor) render() error {
	e.view.Clear()
	err := e.Mode.Render(e.view)
	if err != nil {
		return fmt.Errorf("error on rendering: %w", err)
	}
	size := e.view.Size
	for line := 0; line < size.L; line++ {
		for col := 0; col < size.C; col++ {
			rune := e.view.Content[line][col]
			selection := e.view.Selection[line][col]
			style := tcell.StyleDefault.Reverse(selection)
			e.screen.SetContent(col, line, rune, nil, style)
		}
	}
	status := []rune(e.view.Status)
	prompt := []rune(e.view.Prompt)
	line := size.L
	for col := 0; col < size.C; col++ {
		rune := ' '
		style := tcell.StyleDefault.Background(colors[e.view.Color])
		if col < len(status) {
			rune = status[col]
		}
		if col >= len(status)+1 && col < len(status)+len(prompt)+1 {
			rune = prompt[col-len(status)-1]
			style = style.Reverse(true)
		}
		e.screen.SetContent(col, line, rune, nil, style)
	}
	switch e.view.Cursor {
	case CursorNone:
		e.screen.HideCursor()
	case CursorContent:
		e.screen.ShowCursor(e.view.Position.C, e.view.Position.L)
	case CursorPrompt:
		e.screen.ShowCursor(len(status)+1+len(prompt), line)
	}
	e.screen.Show()
	return nil
}
