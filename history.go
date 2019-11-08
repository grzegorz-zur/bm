package main

const (
	HistorySize = 1024
)

type History struct {
	Records
	bottom  int
	current int
	top     int
}

type Records []Record

type Record struct {
	Lines
	Position
}

func (history *History) Archive(lines Lines, position Position) {
	if history.Records == nil {
		history.Records = make(Records, HistorySize)
	} else {
		history.top = wrap(history.top, HistorySize, 1, Forward)
		if history.top == history.bottom {
			history.bottom = wrap(history.bottom, HistorySize, 1, Forward)
		}
	}
	history.current = history.top
	record := Record{
		Lines:    lines,
		Position: position,
	}
	history.Records[history.current] = record
}

func (history *History) Switch(dir Direction) (lines Lines, position Position) {
	if dir == Backward && history.current != history.bottom ||
		dir == Forward && history.current != history.top {
		history.current = wrap(history.current, HistorySize, 1, dir)
	}
	record := history.Records[history.current]
	lines, position = record.Lines, record.Position
	return
}
