package main

func wrap(position, length, step int, direction Direction) int {
	if length == 0 {
		return 0
	}
	return ((position+step*int(direction))%length + length) % length
}
