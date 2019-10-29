package bm

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
			pos, ok = file.advance(pos, dir)
			if !ok {
				return file.Position
			}
			if file.atWord(pos) {
				return pos
			}
		}
	}
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

func (file File) runeAt(pos Position) rune {
	return file.Lines[pos.Line][pos.Col]
}

func (file File) atText(pos Position) bool {
	if pos.Line >= len(file.Lines) {
		return false
	}
	l := file.Lines[pos.Line]
	if pos.Col >= len(l) {
		return false
	}
	return true
}

func (file File) advance(pos Position, dir Direction) (next Position, ok bool) {
	next = pos
	if dir == Backward && pos.Line == 0 && pos.Col == 0 {
		return
	}
	if dir == Forward && (pos.Line >= len(file.Lines) ||
		pos.Line == len(file.Lines)-1 && pos.Col >= len(file.Lines[pos.Line])) {
		return
	}
	line := pos.Line
	col := pos.Col + dir.Value()
	if col < 0 {
		line--
		col = len(file.Lines[line]) - 1
	}
	if col >= len(file.Lines[line]) {
		line++
		col = 0
	}
	if line < 0 || line >= len(file.Lines) {
		return
	}
	ok = true
	next = Position{Line: line, Col: col}
	return
}
