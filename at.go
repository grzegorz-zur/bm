package main

import (
	"unicode"
)

// AtFileStart checks if cursor is at file start.
func (file *File) AtFileStart() bool {
	return file.location == 0
}

// AtFileEnd checks if cursor is at file end.
func (file *File) AtFileEnd() bool {
	return file.location == len(file.content)
}

// AtLineStart checks if cursor is at line start.
func (file *File) AtLineStart() bool {
	if file.AtFileStart() {
		return true
	}
	last, _ := file.last()
	return last == EOL
}

// AtLineEnd checks if cursor is at line end.
func (file *File) AtLineEnd() bool {
	if file.AtFileEnd() {
		return true
	}
	current, _ := file.current()
	return current == EOL
}

// AtWord checks if cursor is at word.
func (file *File) AtWord() bool {
	if file.AtFileStart() || file.AtFileEnd() {
		return true
	}
	last, _ := file.last()
	current, _ := file.current()
	return !(unicode.IsLetter(last) || unicode.IsDigit(last)) &&
		(unicode.IsLetter(current) || unicode.IsDigit(current))
}

// AtParagraph checks if cursor is at paragraph.
func (file *File) AtParagraph() bool {
	if file.AtFileStart() || file.AtFileEnd() {
		return true
	}
	last, size := file.last()
	location := file.location - size
	if location == 0 && last == EOL {
		return true
	}
	previous, _ := file.previous(location)
	return previous == EOL && last == EOL
}
