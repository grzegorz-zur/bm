package main

import (
	"unicode"
)

// Motion represents cursor movement in the file.
type Motion func(File) Position

// Left moves cursor to the left.
func (f File) Left() Position {
	p := f.pos
	if p.C > 0 {
		p.C--
	}
	return p
}

// Right moves cursor to the right.
func (f File) Right() Position {
	p := f.pos
	p.C++
	return p
}

// Up moves cursor line up.
func (f File) Up() Position {
	p := f.pos
	if p.L > 0 {
		p.L--
	}
	return p
}

// Down moves cursor line down.
func (f File) Down() Position {
	p := f.pos
	p.L++
	return p
}

// Word moves cursor to the next or previous word.
func Word(d Direction) Motion {
	return func(f File) Position {
		p := f.pos
		for {
			var ok bool
			p, ok = f.nextRune(p, d)
			if !ok {
				return f.pos
			}
			if f.atWord(p) {
				return p
			}
		}
	}
}

// Paragraph moves cursor to the next or previous paragraph.
func Paragraph(d Direction) Motion {
	return func(f File) Position {
		p := f.pos
		for {
			var ok bool
			p, ok = f.nextLine(p, d)
			if !ok {
				return f.pos
			}
			p.C = 0
			if f.atParagraph(p) {
				return p
			}
		}
	}
}

func (ls Lines) nextRune(p Position, d Direction) (Position, bool) {
	if p.L < len(ls) {
		c := p.C + d.Value()
		if 0 <= c && c < len(ls[p.L]) {
			p.C = c
			return p, true
		} else if d == Backward && c >= 0 && len(ls[p.L]) > 0 {
			p.C = len(ls[p.L]) - 1
			return p, true
		}
	}
	p, ok := ls.nextLine(p, d)
	p.C = 0
	if d == Backward && len(ls) > p.L && len(ls[p.L]) > 0 {
		p.C = len(ls[p.L]) - 1
	}
	return p, ok
}

func (ls Lines) nextLine(p Position, d Direction) (Position, bool) {
	if p.L == 0 && d == Backward {
		return p, false
	}
	if p.L >= len(ls) && d == Forward {
		return p, false
	}
	if p.L >= len(ls) && d == Backward && len(ls) > 0 {
		p.L = len(ls) - 1
		return p, true
	}
	for l := p.L + d.Value(); 0 <= l && l < len(ls); l += d.Value() {
		if len(ls[l]) > 0 {
			p.L = l
			return p, true
		}
	}
	return p, false
}

func (ls Lines) runeAt(p Position) rune {
	return ls[p.L][p.C]
}

func (ls Lines) atWord(p Position) bool {
	if !ls.AtText(p) {
		return false
	}
	r := ls.runeAt(p)
	if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
		return false
	}
	if p.C == 0 {
		return true
	}
	p.C--
	r = ls.runeAt(p)
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	return true
}

func (ls Lines) atParagraph(p Position) bool {
	if !ls.AtText(p) {
		return false
	}
	if len(ls[p.L]) == 0 {
		return false
	}
	if p.L == 0 {
		return true
	}
	p.L--
	if len(ls[p.L]) == 0 {
		return true
	}
	return false
}
