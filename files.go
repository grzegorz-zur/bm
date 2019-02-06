package bm

import "github.com/pkg/errors"

type Files struct {
	*File
	list []*File
}

func (files *Files) Empty() bool {
	return len(files.list) == 0
}

func (files *Files) Open(path string) (err error) {
	index, found := files.find(path)
	if found {
		files.switchFile(index)
		return
	}
	file, err := Open(path)
	if err != nil {
		err = errors.Wrapf(err, "error opening file %s", file.Path)
		return
	}
	index = files.add(&file)
	files.switchFile(index)
	return
}

func (files *Files) SwitchFile(dir Direction) {
	if files.Empty() {
		return
	}
	index := files.current()
	index = wrap(index, len(files.list), 1, dir)
	files.switchFile(index)
	files.ReloadIfModified()
}

func (files *Files) WriteAll() (err error) {
	if files.Empty() {
		return
	}
	for _, file := range files.list {
		err = file.Write()
		if err != nil {
			err = errors.Wrapf(err, "error writing file %s", file.Path)
		}
	}
	return
}

func (files *Files) Close() {
	if files.Empty() {
		return
	}
	index := files.current()
	files.remove(index)
	index = wrap(index, len(files.list), 0, Forward)
	files.switchFile(index)
}

func (files *Files) add(file *File) (index int) {
	files.list = append(files.list, file)
	index = len(files.list) - 1
	return
}

func (files *Files) remove(index int) {
	list := make([]*File, 0, len(files.list)-1)
	list = append(list, files.list[:index]...)
	list = append(list, files.list[index+1:]...)
	files.list = list
}

func (files *Files) switchFile(index int) {
	if len(files.list) != 0 {
		files.File = files.list[index]
	} else {
		files.File = nil
	}
}

func (files *Files) find(path string) (index int, found bool) {
	for i, file := range files.list {
		if file.Path == path {
			return i, true
		}
	}
	return
}

func (files *Files) current() (index int) {
	for i, file := range files.list {
		if files.File == file {
			return i
		}
	}
	return
}
