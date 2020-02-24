package main

// Mode represents a mode of the editor.
//
// Mode handles input and produces content.
type Mode interface {
	Show() error
	Hide() error
	Key(Key) error
	Rune(rune) error
	Render(*View) error
}
