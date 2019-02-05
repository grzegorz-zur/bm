package bm

func wrap(pos, length, step int, dir Direction) int {
	if length == 0 {
		return 0
	}
	return ((pos+step*dir.Value())%length + length) % length
}
