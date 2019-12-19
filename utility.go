package main

func wrap(pos, len, step int, d Direction) int {
	if len == 0 {
		return 0
	}
	return ((pos+step*d.Value())%len + len) % len
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
