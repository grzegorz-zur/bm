package main

func wrap(pos, len, step int, d Direction) int {
	if len == 0 {
		return 0
	}
	return ((pos+step*d.Value())%len + len) % len
}
