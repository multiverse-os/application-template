package mumble

import (
	"fmt"
	"os"
	"path/filepath"
)

// NOTE
// The goal is to be a very simple filesystem interface, simplifying interaction
// with the filesystem by abstracting a thin as possible layer making code
// expressive as possible. To be successful this file must stay small and not
// complex at all; but also should make working with filesystems very natural,
// and even have validation for security.

type Path string
type File Path
type Directory Path

func (self Path) String() string { return string(self) }

//func (self Directory) Path() string { return self.path.String() }

func (self Path) Directory(directory string) Path {
	return Path(fmt.Sprintf("%s/%s/", self.String(), directory))
}

func (self Directory) Name() string { return filepath.Base(Path(self).String()) }

func (self File) Name() string {
	return filepath.Base(Path(self).String())
}

func (self File) Basename() string {
	return self.Name()[0:(len(self.Name()) - len(self.Extension()))]
}

func (self File) Extension() string { return filepath.Ext(Path(self).String()) }

func (self Directory) File(filename string) Path {
	return Path(fmt.Sprintf("%s/%s", Path(self).String(), filename))
}

// NOTE: Create directories if they don't exist, or simply create the
// directory, so we can have a single create for either file or directory.
func (self Path) Create(permissions os.FileMode) (*os.File, error) {
	baseDirectory := filepath.Dir(self.String())
	switch baseDirectory {
	case ".", "..", "/":
		return nil, fmt.Errorf("error: already exists")
	default:
		os.MkdirAll(baseDirectory, os.ModePerm)
	}
	return os.OpenFile(self.String(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, permissions)
}

func (self Path) Move(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	Path(path).Create(info.Mode())
	return self.Remove()
}

func (self Path) Rename(path string) error {
	baseDirectory := filepath.Dir(path)
	os.MkdirAll(baseDirectory, os.ModePerm)
	return os.Rename(self.String(), path)
}

func (self Path) Remove() error {
	return os.RemoveAll(self.String())
}

// Validation /////////////////////////////////////////////////////////////////
func (self Path) Clean() Path {
	path := filepath.Clean(self.String())
	if filepath.IsAbs(path) {
		return Path(path)
	} else {
		path, _ = filepath.Abs(path)
		return Path(path)
	}
}

// Checking ///////////////////////////////////////////////////////////////////
func (self Path) Exists() bool {
	_, err := os.Stat(self.String())
	return !os.IsNotExist(err)
}

func (self Path) IsDirectory() bool {
	info, _ := os.Stat(self.String())
	return info.IsDir()
}

func (self Path) IsFile() bool {
	info, _ := os.Stat(self.String())
	return !info.IsDir()
}
