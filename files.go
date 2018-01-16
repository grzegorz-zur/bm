package bm

type Files struct {
	list    []*File
	current int
}

func (files *Files) Change(op Change) {
	*(files.list[files.current]) = op(*(files.list[files.current]))
}

func (files *Files) Move(m Move) {
	files.list[files.current].Move(m)
}

func (files *Files) New() {
	file := NewFile()
	files.current = files.add(&file)
}

func (files *Files) Open(path string) (err error) {
	file, err := Read(path)
	if err != nil {
		return
	}
	files.current = files.add(&file)
	return
}

func (files *Files) Current() *File {
	if len(files.list) == 0 {
		files.New()
	}
	return files.list[files.current]
}

func (files *Files) Next(dir Direction) {
	files.current = wrap(files.current, len(files.list), 1, dir)
}

func (files *Files) Write() (err error) {
	err = files.list[files.current].Write()
	return
}

func (files *Files) Close() {
	files.remove(files.current)
	files.current = wrap(files.current, len(files.list), 0, Forward)
}

func (files *Files) add(file *File) (position int) {
	files.list = append(files.list, file)
	position = len(files.list) - 1
	return
}

func (files *Files) remove(position int) {
	list := make([]*File, len(files.list)-1)
	list = append(list, files.list[:position]...)
	list = append(list, files.list[position+1:]...)
	files.list = list
	return
}

func wrap(position, length, step int, dir Direction) int {
	if length == 0 {
		return 0
	}
	return ((position+step*dir.Value())%length + length) % length
}
