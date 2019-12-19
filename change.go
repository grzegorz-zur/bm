package main

// Change represents a file content modification.
type Change func(File) File

// DeleteRune deletes a character at the current position.
func (f File) DeleteRune() File {
	f.Lines = f.Lines.DeleteRune(f.pos)
	return f
}

// DeletePreviousRune deletes a character left of the current position.
func (f File) DeletePreviousRune() File {
	if f.pos.C == 0 {
		return f
	}
	f.Lines = f.Lines.DeletePreviousRune(f.pos)
	f.pos.C--
	return f
}

// InsertRune inserts a character at the current position.
func InsertRune(r rune) Change {
	return func(f File) File {
		f.Lines = f.Lines.InsertRune(f.pos, r)
		f.pos.C++
		return f
	}
}

// DeleteLine deletes a line at the current position.
func (f File) DeleteLine() File {
	f.Lines = f.Lines.DeleteLine(f.pos.L)
	return f
}

// Split splites a line at the current position.
func (f File) Split() File {
	f.Lines = f.Lines.Split(f.pos)
	f.pos.L++
	f.pos.C = 0
	return f
}

// PasteBlock inserts multiple lines at the current position as a block.
func PasteBlock(ls Lines) Change {
	return func(f File) File {
		f.Lines = f.Lines.InsertBlock(f.pos, ls)
		return f
	}
}

// PasteInline inserts multiple lines at the current position inline.
func PasteInline(ls Lines) Change {
	return func(f File) File {
		f.Lines = f.Lines.InsertInline(f.pos, ls)
		return f
	}
}

// Delete deletes content between positions.
func (f File) Delete() File {
	s, e := Sort(f.sel, f.pos)
	f.Lines = f.Lines.Delete(s, e)
	return f
}
