package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Directory struct {
	Path        Path
	Directories []*Directory
	Files       []*File
}

func (self *Directory) NewDirectory(name string) error {
	path := filepath.Join(self.Path, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.FileMode(0770))
		if err != nil {
			return err
		} else {
			self.Directories = append(self.Directories, &Directory{
				Path:        path,
				Directories: []*Directory{},
				Files:       []*File{},
			})
		}
	}
}

// TODO: Should maybe overrite?
func (self *Directory) NewFile(name string, data []byte) error {
	path := filepath.Join(self.Path, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := ioutil.WriteFile(path, data, 0644)
		if err != nil {
			return err
		} else {
			self.Files = append(self.Files, &File{
				Path: path,
				Data: data,
			})
		}
	}
}
