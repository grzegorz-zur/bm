package bm

type Files struct {
	*File
	list []*File
}

func (files *Files) Empty() bool {
	return files.File == nil
}

func (files *Files) Open(base, path string) (err error) {
	position, found := files.find(path)
	if found {
		files.switchFile(position)
		return
	}
	file, err := Open(base, path)
	if err != nil {
		return
	}
	position = files.add(&file)
	files.switchFile(position)
	return
}

func (files *Files) SwitchFile(dir Direction) {
	if files.Empty() {
		return
	}
	position := files.current()
	position = wrap(position, len(files.list), 1, dir)
	files.switchFile(position)
}

func (files *Files) WriteAll(base string) (err error) {
	for _, file := range files.list {
		err = file.Write(base)
		if err != nil {
			return
		}
	}
	return
}

func (files *Files) Close() {
	if files.Empty() {
		return
	}
	position := files.current()
	files.remove(position)
	position = wrap(position, len(files.list), 0, Forward)
	files.switchFile(position)
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

func (files *Files) switchFile(position int) {
	if len(files.list) != 0 {
		files.File = files.list[position]
	} else {
		files.File = nil
	}
}

func (files *Files) find(path string) (position int, found bool) {
	for i, file := range files.list {
		if file.Path == path {
			return i, true
		}
	}
	return
}

func (files *Files) current() (position int) {
	for position, file := range files.list {
		if files.File == file {
			return position
		}
	}
	panic("failed to find current file: " + files.Path)
}
