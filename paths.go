package main

import (
	"os"
	"strings"
)

// Paths represents list of file paths.
type Paths []string

// Len returns length of paths.
func (paths Paths) Len() int {
	return len(paths)
}

// Less compares elements.
func (paths Paths) Less(i, j int) bool {
	pi := paths[i]
	pj := paths[j]
	sep := string(os.PathSeparator)
	si := strings.Count(pi, sep)
	sj := strings.Count(pj, sep)
	if si == sj {
		return pi < pj
	}
	return si < sj
}

// Swap exchanges elements.
func (paths Paths) Swap(i, j int) {
	paths[i], paths[j] = paths[j], paths[i]
}
