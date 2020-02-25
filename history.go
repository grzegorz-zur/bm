package main

const (
	// HistorySize is length of the history buffer.
	HistorySize = 1024
)

// History holds latest records in a ring buffer.
type History struct {
	records Records
	bottom  int
	current int
	top     int
}

// Records is a list of history records.
type Records []Record

// Record holds content and position.
type Record struct {
	content  string
	location int
}

// Archive saves content and position in history and sets position to the recent entry.
func (history *History) Archive(content string, location int) {
	if history.records == nil {
		history.records = make(Records, HistorySize)
	} else {
		history.top = wrap(history.top, HistorySize, 1, Forward)
		if history.top == history.bottom {
			history.bottom = wrap(history.bottom, HistorySize, 1, Forward)
		}
	}
	history.current = history.top
	record := Record{
		content:  content,
		location: location,
	}
	history.records[history.current] = record
}

// Switch retrieves content and position from history and moves current position.
func (history *History) Switch(direction Direction) (content string, location int) {
	if direction == Backward && history.current != history.bottom ||
		direction == Forward && history.current != history.top {
		history.current = wrap(history.current, HistorySize, 1, direction)
	}
	record := history.records[history.current]
	return record.content, record.location
}
