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
	content   *Content
	buffer    Lines
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
		content:   NewContent(Size{}),
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
	e.buffer = e.Selection()
}

// PasteBlock pastes buffer to current file as a block.
func (e *Editor) PasteBlock() {
	e.Change(PasteBlock(e.buffer))
}

// PasteInline pastes buffer to current file inline.
func (e *Editor) PasteInline() {
	e.Change(PasteInline(e.buffer))
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

func (e *Editor) handle(ev tcell.Event) error {
	switch evt := ev.(type) {
	case *tcell.EventKey:
		if evt.Key() == tcell.KeyRune {
			return e.Mode.Rune(evt.Rune())
		}
		k, ok := keymap[evt.Key()]
		if ok {
			return e.Mode.Key(k)
		}
	case *tcell.EventResize:
		w, h := evt.Size()
		if h > 0 {
			h--
		}
		s := Size{h, w}
		e.content = NewContent(s)
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
	e.content.Clear()
	err := e.Mode.Render(e.content)
	if err != nil {
		return fmt.Errorf("error on rendering: %w", err)
	}
	s := e.content.Size
	for l := 0; l < s.L; l++ {
		for c := 0; c < s.C; c++ {
			r := e.content.Runes[l][c]
			m := e.content.Marks[l][c]
			stl := tcell.StyleDefault.Reverse(m)
			e.screen.SetContent(c, l, r, nil, stl)
		}
	}
	rs := []rune(e.content.Status)
	rp := []rune(e.content.Prompt)
	l := s.L
	for c := 0; c < s.C; c++ {
		r := ' '
		stl := tcell.StyleDefault.Background(colors[e.content.Color])
		if c < len(rs) {
			r = rs[c]
		}
		if c >= len(rs)+1 && c < len(rs)+len(rp)+1 {
			r = rp[c-len(rs)-1]
			stl = stl.Reverse(true)
		}
		e.screen.SetContent(c, l, r, nil, stl)
	}
	switch e.content.Cursor {
	case CursorNone:
		e.screen.HideCursor()
	case CursorContent:
		e.screen.ShowCursor(e.content.Position.C, e.content.Position.L)
	case CursorPrompt:
		e.screen.ShowCursor(len(rs)+1+len(rp), l)
	}
	e.screen.Show()
	return nil
}
