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
	screenCreate ScreenCreate
	screen       tcell.Screen
	view         *View
	content      string
	events       chan tcell.Event
	check        chan struct{}
	pause        chan struct{}
	unpause      chan struct{}
	quit         chan struct{}
	done         chan struct{}
}

// ScreenCreate creates a new screen.
type ScreenCreate func() (tcell.Screen, error)

// New creates new editor.
func New(screenCreate ScreenCreate, paths []string) *Editor {
	editor := &Editor{
		screenCreate: screenCreate,
		view:         NewView(Size{}),
		events:       make(chan tcell.Event),
		check:        make(chan struct{}),
		pause:        make(chan struct{}, 1),
		unpause:      make(chan struct{}),
		quit:         make(chan struct{}, 1),
		done:         make(chan struct{}),
	}
	editor.Modes.Command = &Command{editor: editor}
	editor.Modes.Input = &Input{editor: editor}
	editor.Modes.Select = &Select{editor: editor}
	editor.Modes.Switch = &Switch{editor: editor}
	for _, path := range paths {
		editor.Open(path)
	}
	editor.SwitchFile(Forward)
	editor.SwitchMode(editor.Command)
	return editor
}

// Start starts the editor.
func (editor *Editor) Start() (err error) {
	editor.screen, err = editor.screenCreate()
	if err != nil {
		return fmt.Errorf("error creating screen: %w", err)
	}
	err = editor.screen.Init()
	if err != nil {
		return fmt.Errorf("error initializing screen: %w", err)
	}
	go editor.signals()
	go editor.listen()
	go editor.tick()
	go editor.run()
	return nil
}

// Check signals file modification check.
func (editor *Editor) Check() {
	editor.check <- struct{}{}
}

// Pause pauses the editor.
func (editor *Editor) Pause() {
	editor.pause <- struct{}{}
}

// Quit quits the editor.
func (editor *Editor) Quit() {
	editor.quit <- struct{}{}
}

// Wait waits for the editor to finish.
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

// Copy copies selection to buffer.
func (editor *Editor) Copy() {
	editor.content = editor.File.Copy()
	editor.SwitchMode(editor.Command)
}

// Cut cuts selection to buffer.
func (editor *Editor) Cut() {
	editor.content = editor.File.Cut()
	editor.SwitchMode(editor.Command)
}

// Paste pastes buffer.
func (editor *Editor) Paste() {
	editor.Insert(editor.content)
}

// LineAbove starts line above current line.
func (editor *Editor) LineAbove() {
	editor.MoveLineStart()
	editor.Insert(string(EOL))
	editor.MoveUp()
	editor.SwitchMode(editor.Input)
}

// LineBelow starts line below current line.
func (editor *Editor) LineBelow() {
	editor.MoveLineEnd()
	editor.Insert(string(EOL))
	editor.SwitchMode(editor.Input)
}

func (editor *Editor) listen() {
	for {
		editor.events <- editor.screen.PollEvent()
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
	render := true
	for {
		if render {
			err := editor.render()
			editor.report(err)
		}
		select {
		case event := <-editor.events:
			err := editor.handle(event)
			editor.report(err)
			render = true
		case <-editor.check:
			if editor.Empty() {
				render = false
			} else {
				read, err := editor.Read(false)
				render = read
				editor.report(err)
			}
		case <-editor.pause:
			err := editor.background()
			editor.report(err)
		case <-editor.quit:
			editor.close()
			return
		}
	}
}

func (editor *Editor) handle(event tcell.Event) error {
	switch tevent := event.(type) {
	case *tcell.EventKey:
		if tevent.Key() == tcell.KeyRune {
			return editor.Rune(tevent.Rune())
		}
		key, ok := keymap[tevent.Key()]
		if ok {
			return editor.Key(key)
		}
	case *tcell.EventResize:
		log.Println(tevent)
		width, height := tevent.Size()
		if height > 0 {
			height--
		}
		size := Size{height, width}
		editor.view = NewView(size)
	}
	return nil
}

func (editor *Editor) report(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (editor *Editor) background() error {
	editor.screen.Fini()
	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("error pausing process %+v: %w", process, err)
	}
	process.Signal(syscall.SIGSTOP)
	<-editor.unpause
	editor.screen, err = editor.screenCreate()
	if err != nil {
		return fmt.Errorf("error creating screen: %w", err)
	}
	err = editor.screen.Init()
	if err != nil {
		return fmt.Errorf("error initializing screen: %w", err)
	}
	return nil
}

func (editor *Editor) close() {
	editor.screen.Fini()
}

func (editor *Editor) render() error {
	editor.view.Clear()
	err := editor.Mode.Render(editor.view)
	if err != nil {
		return fmt.Errorf("error on rendering: %w", err)
	}
	size := editor.view.Size
	for line := 0; line < size.Lines; line++ {
		for column := 0; column < size.Columns; column++ {
			rune := editor.view.Content[line][column]
			selection := editor.view.Selection[line][column]
			style := tcell.StyleDefault.Reverse(selection)
			editor.screen.SetContent(column, line, rune, nil, style)
		}
	}
	status := []rune(editor.view.Status)
	prompt := []rune(editor.view.Prompt)
	line := size.Lines
	for column := 0; column < size.Columns; column++ {
		rune := ' '
		style := tcell.StyleDefault.Background(colors[editor.view.Color])
		if column < len(status) {
			rune = status[column]
		}
		if column >= len(status)+1 && column < len(status)+len(prompt)+1 {
			rune = prompt[column-len(status)-1]
			style = style.Reverse(true)
		}
		editor.screen.SetContent(column, line, rune, nil, style)
	}
	switch editor.view.Cursor {
	case CursorNone:
		editor.screen.HideCursor()
	case CursorContent:
		editor.screen.ShowCursor(editor.view.Position.Column, editor.view.Position.Line)
	case CursorPrompt:
		editor.screen.ShowCursor(len(status)+1+len(prompt), line)
	}
	editor.screen.Show()
	return nil
}
