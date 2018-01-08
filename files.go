package bm

type Files []*File

func (fs Files) Add(f *File) Files {
	return append(fs, f)
}

func (fs Files) Next(f *File, d Direction) *File {
	i := fs.find(f)
	k := int(d)
	n := len(fs)
	j := ((i+k)%n + n) % n
	return fs[j]
}

func (fs Files) find(f *File) int {
	for i := range fs {
		if fs[i] == f {
			return i
		}
	}
	panic("file " + f.Path + " not on the files")
}
