package main

import (
	"unicode"
)

type Move func(file File) (pos Position)

func (file File) Left() (pos Position) {
	pos = file.Position
	if pos.Col > 0 {
		pos.Col--
	}
	return
}

func (file File) Right() (pos Position) {
	pos = file.Position
	pos.Col++
	return
}

func (file File) Up() (pos Position) {
	pos = file.Position
	if pos.Line > 0 {
		pos.Line--
	}
	return
}

func (file File) Down() (pos Position) {
	pos = file.Position
	pos.Line++
	return
}

func (file File) Word(dir Direction) Move {
	return func(file File) Position {
		pos := file.Position
		for {
			var ok bool
			pos, ok = file.nextRune(pos, dir)
			if !ok {
				return file.Position
			}
			if file.atWord(pos) {
				return pos
			}
		}
	}
}

func (file File) nextRune(pos Position, dir Direction) (Position, bool) {
	l := pos.Line
	if l < len(file.Lines) {
		c := pos.Col + dir.Value()
		if 0 <= c && c < len(file.Lines[l]) {
			return Position{l, c}, true
		}
	}
	return file.nextLine(pos, dir)
}

func (file File) nextLine(pos Position, dir Direction) (Position, bool) {
	if pos.Line == 0 && dir == Backward {
		return pos, false
	}
	if pos.Line >= len(file.Lines) && dir == Forward {
		return pos, false
	}
	for l := pos.Line + dir.Value(); 0 <= l && l < len(file.Lines); l += dir.Value() {
		if len(file.Lines[l]) > 0 {
			c := 0
			if dir == Backward {
				c = len(file.Lines[l]) - 1
			}
			return Position{l, c}, true
		}
	}
	return pos, false
}

func (file File) runeAt(pos Position) rune {
	return file.Lines[pos.Line][pos.Col]
}

func (file File) atWord(pos Position) bool {
	if !file.atText(pos) {
		return false
	}
	r := file.runeAt(pos)
	if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
		return false
	}
	if pos.Col == 0 {
		return true
	}
	pos.Col--
	r = file.runeAt(pos)
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	return true
}

func (file File) atText(pos Position) bool {
	if pos.Line >= len(file.Lines) {
		return false
	}
	if pos.Col >= len(file.Lines[pos.Line]) {
		return false
	}
	return true
}
