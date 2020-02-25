package main

// Insert inserts content to the file.
func (file *File) Insert(content string) {
	file.content = file.content[:file.location] + content + file.content[file.location:]
	file.location += len(content)
	file.changed = true
	file.Archive()
}

// Backspace deletes rune on the left.
func (file *File) Backspace() {
	if file.AtFileStart() {
		return
	}
	file.MoveLeft()
	file.Delete()
}

// Delete removes current rune.
func (file *File) Delete() {
	if file.AtFileEnd() {
		return
	}
	_, size := file.current()
	from := file.location
	to := from + size
	file.Remove(from, to)
}

// DeleteLine removes current line.
func (file *File) DeleteLine() {
	location := file.location
	file.MoveLineStart()
	from := file.location
	file.MoveLineEnd()
	file.MoveRight()
	to := file.location
	file.location = location
	file.Remove(from, to)
}

// Remove deletes content from the file.
func (file *File) Remove(from, to int) {
	file.content = file.content[:from] + file.content[to:]
	if file.location > from {
		file.location = from
	}
	file.changed = true
	file.Archive()
}
