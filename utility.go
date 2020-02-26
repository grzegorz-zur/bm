package main

func wrap(position, length, step int, direction Direction) int {
	return ((position+step*int(direction))%length + length) % length
}

func visible(rune rune) rune {
	switch {
	case rune < 16:
		return rune + '␀'
	case rune == 255:
		return '␡'
	}
	return rune
}
