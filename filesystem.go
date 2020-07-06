package mumble

import (
	"fmt"
	"os"
)

type Directory struct {
	Name string
	path Path
}

type Path string

func (self Path) String() string { return string(self) }

func (self Directory) Path() string { return self.path.String() }

func (self Directory) Subdirectory(directory string) string {
	return fmt.Sprintf("%s/%s/", self.path.String(), directory)
}

func (self Directory) File(filename string) Path {
	return Path(fmt.Sprintf("%s/%s", self.path.String(), filename))
}

func (self Path) Create(permissions os.FileMode) (*os.File, error) {
	return os.OpenFile(self.String(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, permissions)
}
