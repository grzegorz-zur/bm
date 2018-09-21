package bm

func wrap(position, length, step int, dir Direction) int {
	if length == 0 {
		return 0
	}
	return ((position+step*dir.Value())%length + length) % length
}
