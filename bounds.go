package bm

type Bounds struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

func (b Bounds) Size() Size {
	return Size{
		Lines: b.Bottom - b.Top,
		Cols:  b.Right - b.Left,
	}
}
