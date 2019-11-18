package main

const (
	// HistorySize is length of the history buffer.
	HistorySize = 1024
)

// History holds latest records in a ring buffer.
type History struct {
	Records
	bottom  int
	current int
	top     int
}

// Records is a list of history records.
type Records []Record

// Record holds content and position.
type Record struct {
	Lines
	Position
}

// Archive saves content and position in history and sets position to the recent entry.
func (h *History) Archive(ls Lines, p Position) {
	if h.Records == nil {
		h.Records = make(Records, HistorySize)
	} else {
		h.top = wrap(h.top, HistorySize, 1, Forward)
		if h.top == h.bottom {
			h.bottom = wrap(h.bottom, HistorySize, 1, Forward)
		}
	}
	h.current = h.top
	r := Record{
		Lines:    ls,
		Position: p,
	}
	h.Records[h.current] = r
}

// Switch retrieves content and position from history and moves current position.
func (h *History) Switch(d Direction) (Lines, Position) {
	if d == Backward && h.current != h.bottom ||
		d == Forward && h.current != h.top {
		h.current = wrap(h.current, HistorySize, 1, d)
	}
	r := h.Records[h.current]
	return r.Lines, r.Position
}
