package bm

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

type Content struct {
	Lines  []string
	Offset Position
	Cursor Position
}

type Position struct {
	Column int
	Line   int
}

func (content *Content) Display(minX, minY, maxX, maxY int) (x, y int, err error) {
	for y := minY; y <= maxY; y++ {
		oy := content.Offset.Line + y
		if oy >= len(content.Lines) {
			break
		}
		line := content.Lines[oy]
		runes := []rune(line)
		for x := minX; x <= maxX; x++ {
			ox := content.Offset.Line + x
			if ox >= len(runes) {
				break
			}
			r := runes[ox]
			tb.SetCell(x, y, r, tb.ColorDefault, tb.ColorDefault)
		}
	}
	x, y = content.Cursor.Column, content.Cursor.Line
	return
}

func (content *Content) Key(event tb.Event) {
	text := fmt.Sprintf("%+v", event)
	content.Lines = append(content.Lines, text)
	content.Cursor.Line = len(content.Lines)
	return
}
