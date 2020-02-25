package main

// MoveLeft moves cursor to the left.
func (file *File) MoveLeft() {
	if file.AtFileStart() {
		return
	}
	_, size := file.last()
	file.location -= size
}

// MoveRight moves cursor to the right.
func (file *File) MoveRight() {
	if file.AtFileEnd() {
		return
	}
	_, size := file.current()
	file.location += size
}

// MoveUp moves cursor up.
func (file *File) MoveUp() {
	location := file.location
	file.MoveLineStart()
	columns := location - file.location
	file.MoveLeft()
	file.MoveLineStart()
	for column := 0; column < columns && !file.AtLineEnd(); column++ {
		file.MoveRight()
	}
}

// MoveDown moves cursor down.
func (file *File) MoveDown() {
	location := file.location
	file.MoveLineStart()
	columns := location - file.location
	file.MoveLineEnd()
	file.MoveRight()
	for column := 0; column < columns && !file.AtLineEnd(); column++ {
		file.MoveRight()
	}
}

// MoveLineStart moves cursor to the start of line.
func (file *File) MoveLineStart() {
	for !file.AtLineStart() {
		file.MoveLeft()
	}
}

// MoveLineEnd moves cursor to the end of line.
func (file *File) MoveLineEnd() {
	for !file.AtLineEnd() {
		file.MoveRight()
	}
}

// MoveFileStart moves cursor to the start of file.
func (file *File) MoveFileStart() {
	file.location = 0
}

// MoveFileEnd moves cursor to the end of file.
func (file *File) MoveFileEnd() {
	file.location = len(file.content)
}

// MoveWordPrevious moves cursor to previous word.
func (file *File) MoveWordPrevious() {
	file.MoveLeft()
	for !file.AtWord() {
		file.MoveLeft()
	}
}

// MoveWordNext moves cursor to the next word.
func (file *File) MoveWordNext() {
	file.MoveRight()
	for !file.AtWord() {
		file.MoveRight()
	}
}

// MoveParagraphPrevious moves cursor to the previous paragraph.
func (file *File) MoveParagraphPrevious() {
	file.MoveLeft()
	for !file.AtParagraph() {
		file.MoveLeft()
	}
}

// MoveParagraphNext moves to cursor the next paragraph.
func (file *File) MoveParagraphNext() {
	file.MoveRight()
	for !file.AtParagraph() {
		file.MoveRight()
	}
}
