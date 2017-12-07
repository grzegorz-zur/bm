package bm

type Bounds struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

func (bounds Bounds) Size() Size {
	return Size{
		Lines: bounds.Bottom - bounds.Top,
		Cols:  bounds.Right - bounds.Left,
	}
}

func (bounds Bounds) SplitHorizontal(line int) (top, bottom Bounds) {
	split := 0
	if line > 0 {
		split = bounds.Top + line
	} else {
		split = bounds.Bottom + line
	}
	top = bounds
	top.Bottom = split
	bottom = bounds
	bottom.Top = split + 1
	return
}
