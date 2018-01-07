package bm

type Files []*File

func (fs Files) Add(f *File) Files {
	return append(fs, f)
}

func (fs Files) Next(f *File) *File {
	i := fs.find(f)
	n := (i + 1) % len(fs)
	return fs[n]
}

func (fs Files) find(f *File) int {
	for i := range fs {
		if fs[i] == f {
			return i
		}
	}
	panic("file not on the files")
}
