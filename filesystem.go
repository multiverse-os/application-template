package mumble

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

////////////////////////////////////////////////////////////////////////////////
// NOTE
// The goal is to be a very simple filesystem interface, simplifying interaction
// with the filesystem by abstracting a thin as possible layer making code
// expressive as possible. To be successful this file must stay small and not
// complex at all; but also should make working with filesystems very natural,
// and even have validation for security.
////////////////////////////////////////////////////////////////////////////////

type Path string
type File Path
type Directory Path

func (self Path) String() string { return string(self) }

////////////////////////////////////////////////////////////////////////////////
func (self Path) Directory(directory string) Path {
	return Path(fmt.Sprintf("%s/%s/", self.String(), directory))
}

func (self Path) File(filename string) Path {
	return Path(fmt.Sprintf("%s/%s", Path(self).String(), filename))
}

///////////////////////////////////////////////////////////////////////////////
// NOTE: Lets always clean before we get to these so no error is possible.
func (self Path) Info() os.FileInfo {
	info, _ := os.Stat(self.String())
	return info
}

func (self Path) Permissions() os.FileMode {
	return self.Info().Mode()
}

///////////////////////////////////////////////////////////////////////////////

func (self Directory) Directory(directory string) Path {
	return Path(self).Directory(directory)
}

func (self Directory) File(filename string) Path {
	return Path(fmt.Sprintf("%s/%s", Path(self).String(), filename))
}

func (self Directory) String() string { return string(self) }
func (self Directory) Path() Path     { return Path(self) }
func (self Directory) Name() string   { return filepath.Base(Path(self).String()) }

///////////////////////////////////////////////////////////////////////////////
func (self File) String() string   { return string(self) }
func (self File) Path() Path       { return Path(self) }
func (self File) Name() string     { return filepath.Base(Path(self).String()) }
func (self File) Basename() string { return self.Name()[0:(len(self.Name()) - len(self.Extension()))] }

// TODO: In a more complete solution, we would also use magic sequence and mime;
// but that would have to be segregated into an interdependent submodule or not
// at all.
func (self File) Extension() string { return filepath.Ext(Path(self).String()) }

// BASIC OPERATIONS ///////////////////////////////////////////////////////////

// NOTE: In the case of directories, may list contents?
func (self File) Open() (*os.File, error) {
	if self.Path().Exists() {
		return os.Open(self.String())
	} else {
		return self.Path().Create(self.Path().Permissions())
	}
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
	if info, err := os.Stat(path); err != nil {
		return err
	} else {
		Path(path).Create(info.Mode())
		return self.Remove()
	}
}

func (self Path) Rename(path string) error {
	baseDirectory := filepath.Dir(path)
	os.MkdirAll(baseDirectory, os.ModePerm)
	return os.Rename(self.String(), path)
}

func (self Path) Remove() error { return os.RemoveAll(self.String()) }

// IO /////////////////////////////////////////////////////////////////////////

func (self File) Bytes() ([]byte, error) {
	return ioutil.ReadFile(self.String())
}

func (self File) HeadBytes(headSize int) ([]byte, error) {
	var headBytes []byte
	file, _ := self.Open()
	readBytes, err := io.ReadAtLeast(file, headBytes, headSize)
	if readBytes != headSize {
		return headBytes, fmt.Errorf("error: failed to complete read: read ", readBytes, " out of ", headSize, "bytes")
	} else {
		return headBytes, err
	}
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
