package mumble

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

////////////////////////////////////////////////////////////////////////////////
// NOTE
// The goal is to be a very simple filesystem interface, simplifying interaction
// with the filesystem by abstracting a thin as possible layer making code
// expressive as possible. To be successful this file must stay small and not
// complex at all; but also should make working with filesystems very natural,
// and even have validation for security.
//
//   * **So far this model benefits greatly from avoiding holding locks or mem**
//   longer than absolutely necessary.
//
//   * It features chainable functionality.
//
////////////////////////////////////////////////////////////////////////////////

// NOTE If we can prevent ALL path errors by validating and cleaning input, we
//      can have an interface without errors ouputs, or at least a choke-point
//      where they would occur; leaving the rest of the API simpler.
// If there is an error, it will be of type *PathError.

type Path string

// TODO: Just maybe these should include the os.File (or not)
type Directory Path
type File Path

// TODO: Chown, Chmod, SoftLink, HardLink, Stream Write, Stream Read, Zero-Copy

func (self Path) String() string { return string(self) }

////////////////////////////////////////////////////////////////////////////////
func (self Path) Directory(directory string) Path {
	return Path(fmt.Sprintf("%s/%s/", self.String(), directory))
}

func (self Path) File(filename string) Path {
	return Path(fmt.Sprintf("%s/%s", Path(self).String(), filename))
}

///////////////////////////////////////////////////////////////////////////////

func (self Directory) Directory(directory string) Path {
	return Path(self).Directory(directory)
}

func (self Directory) String() string { return string(self) }
func (self Directory) Path() Path     { return Path(self) }
func (self Directory) Name() string   { return filepath.Base(Path(self).String()) }

func (self Directory) File(filename string) (File, error) {
	if 0 < len(filename) {
		return File(fmt.Sprintf("%s/%s", Path(self).String(), filename)), nil
	} else {
		if self.Path().IsFile() {
			return File(self), nil
		} else {
			return File(self), fmt.Errorf("error: path does not resolve to file")
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
func (self File) Directory() (Directory, error) {
	if self.Path().IsDirectory() {
		return Directory(self), nil
	} else {
		return Directory(self), fmt.Errorf("error: path does not resolve to directory")
	}
}

func (self File) BaseDirectory() Directory {
	if path, err := filepath.Abs(self.String()); err != nil {
		panic(err)
	} else {
		return Directory(filepath.Dir(path))
	}
}

func (self File) String() string   { return string(self) }
func (self File) Path() Path       { return Path(self) }
func (self File) Name() string     { return filepath.Base(Path(self).String()) }
func (self File) Basename() string { return self.Name()[0:(len(self.Name()) - len(self.Extension()))] }

// TODO: In a more complete solution, we would also use magic sequence and mime;
// but that would have to be segregated into an interdependent submodule or not
// at all.
func (self File) Extension() string { return filepath.Ext(Path(self).String()) }

// BASIC OPERATIONS ///////////////////////////////////////////////////////////
// NOTE: Create directories if they don't exist, or simply create the
// directory, so we can have a single create for either file or directory.
func (self Path) Move(path string) error {
	if info, err := os.Stat(path); err != nil {
		return err
	} else {
		self.Create()
		return self.Remove()
	}
}

func (self Path) Rename(path string) error {
	baseDirectory := filepath.Dir(path)
	os.MkdirAll(baseDirectory, os.ModePerm)
	return os.Rename(self.String(), path)
}

func (self Path) Remove() error { return os.RemoveAll(self.String()) }

// INFO / META ////////////////////////////////////////////////////////////////
//type OwnershipType bool
//
//const (
//	User OwnershipType = iota
//	Group
//)

//type FileDescriptor struct {
//	Pointer *uinptr
//	*os.File
//}

// NOTE: Lets always clean before we get to these so no error is possible.
func (self Path) Metadata() os.FileInfo {
	info, err := os.Stat(self.String())
	if err != nil {
		panic(err)
	}
	return info
}

//func (self Path) FileDescriptor() {}

//func (self Path) Owner() User, Group {}
//func (self Path) OwnerIDs() UID, GUID {}

func (self Path) GUID() int {
	if stat, ok := self.Metadata().Sys().(*syscall.Stat_t); ok {
		return int(stat.Gid)
	} else {
		panic(fmt.Errorf("error: failed to obtain guid of: ", self.String()))
	}
}

//func (self Path) UID() int {
//	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
//		return int(stat.Uid)
//	} else {
//		panic(fmt.Errorf("error: failed to obtain uid of: ", self.String()))
//	}
//}

// Perm returns the Unix permission bits in m.
//func (m FileMode) Perm() FileMode {
//	return m & ModePerm
//}
func (self Path) Permissions() os.FileMode {
	return self.Metadata().Mode()
}

// IO /////////////////////////////////////////////////////////////////////////
// File is the real representation of *File.
// The extra level of indirection ensures that no clients of os
// can overwrite this data, which could cause the finalizer
// to close the wrong file descriptor.
// type file struct {
// 	pfd         poll.FD
// 	name        string
// 	dirinfo     *dirInfo // nil unless directory being read
// 	nonblock    bool     // whether we set nonblocking mode
// 	stdoutOrErr bool     // whether this is stdout or stderr
// 	appendMode  bool     // whether file is opened for appending
// }

// NOTE: This very important; unlike the stdlibrary, the create action is
// entirely segregated from read, write and sync. This is simply a create action
// only. If it does not exist, it creates it, with the option of overwriting it,
// or append (which since we are segregating read/write/sync is non-action.

// ORIGINAL STDLIB CREATE DOES
// Create creates or truncates the named file. If the file already exists, it
// is truncated. If the file does not exist, it is created with mode 0666
// (before umask). If successful, methods on the returned File can be used for
// I/O; the associated file descriptor has mode O_RDWR. If there is an error,
// it will be of type *PathError.

// WHATS DIFFERENT
// (1) All paths are validated/cleaned so there is a single choke-point of
// PathError's, meaning there is no error return; all functions only return the
// single value, making them easier to chain together, and keeping error
// handling to a single part of the software.
// (2) No value is returned from create; instead an immediate close is called.
// nothing is held in memory, and no locks are held. Instead path information,
// about the file is passed around, so a fine-grained READONLY, WRITEONLY, SYNC
// style IO can be called specifically, with a time limit, on a speicifc segment
// of the file and immediately closed. Overwrite() will create a file if it
// doesn't exist, and overwrite any existing files; and Create(false) will
// create a file if it does not exist, but not overwrite an existing file.

// This gives the the following API
// 		File("some/path/to/file.yaml").Create().Read()
//    File("some/file.yaml").Overwrite().Write("test")
//
// Closes happen automatically, everything is streaming, zero-copy. We are not
// passing around a *os.File, we are passing around a `type Path string` of the
// location.

// **And perhaps we should be just passing around the FD pointer,
// since that would be a uint and not a string, therefore faster lookups?**

// **Still need to handle NewFile(fd) functionaliy**
// NOTE: All panic(err) are temporary just for debugging, while we push all errors
// to the choke-point of validation and cleaning.

func (self Path) Create() {
	switch {
	case self.IsDirectory():
		Directory(self).Create()
	case self.IsFile():
		File(self).Create()
	default:
		panic(fmt.Errorf("error: unsupported type"))
	}
}

func (self Directory) Create() {
	if self.Path().IsFile() {
		File(self.Path()).Create()
	} else {
		if !self.Path().Exists() {
			err := os.MkdirAll(self.String(), 0700|os.ModeSticky)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (self File) Create() File {
	if self.Path().IsDirectory() {
		Directory(self.Path()).Create()
	} else {
		path, err := filepath.Abs(self.String())
		if err != nil {
			panic(err)
		}
		Directory(filepath.Dir(path)).Create()
		if !self.Path().Exists() {
			file, err := os.OpenFile(self.String(), os.O_CREATE|os.O_WRONLY, 0640|os.ModeSticky)
			if err != nil && file.Close() != nil {
				panic(err)
			}
		}
	}
	return self
}

func (self File) Overwrite() File {
	if self.Path().IsDirectory() {
		Directory(self.Path()).Create()
	} else {
		path, err := filepath.Abs(self.String())
		if err != nil {
			panic(err)
		}
		Directory(filepath.Dir(path)).Create()
		file, err := os.OpenFile(self.String(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640|os.ModeSticky)
		if err != nil && file.Close() != nil {
			panic(err)
		}
	}
	return self
}

// TODO: Maybe ChangePermissions to match the expected Chmod and changeowner
// chown? Right now its Path -> read, and file and directroy are set. Which is
// not exsactly natural
func (self File) Permissions(permissions os.FileMode) File {
	err := os.Chmod(self.String(), permissions)
	if err != nil {
		panic(err)
	}
	return self
}

func (self File) Chmod(permissions os.FileMode) File {
	return self.Permissions(permissions)
}

func (self File) Owner(username string) File {
	u, err := user.Lookup(username)
	var uid int
	if u != nil {
		uid := u.Uid
	} else if err != nil {
		user, idError := user.LookupId(username)
		if idError != nil {
			panic(err)
		} else {
			uid, err = strconv.Atoi(u.Uid)
			if err != nil {
				panic(err)
			}
		}
	}
	os.Chown(self.String(), uid, self.Path().GUID())
	return self
}

func (self File) Group(guid int) File {
	return self
}
func (self File) Chown(uid, guid int) File {
	return self
}

// NOTE: In the case of directories, may list contents?
func (self File) Open() *os.File {
	if !self.Path().Exists() {
		self = self.Create()
	}
	openedFile, err := os.Open(self.String())
	if err != nil {
		panic(err)
	}
	return openedFile
}

// IO: Reads //////////////////////////////////////////////////////////////////

func (self File) Bytes() (output []byte) {
	if self.Path().Exists() {
		output, err := ioutil.ReadFile(self.String())
		if err != nil {
			// TODO: For now we will panic on errors so we can catch any that slip
			// by and squash them or move them downstream to the validation/cleanup
			// chokepoint.
			panic(err)
		}
	}
	return output
}

//func (self File) Offset(offsetSize int) *os.File {
//		self.Open()
//		return se.Seek(6, 0)
//}

func (self File) Limit(limitSize int) ([]byte, error) {
	var headBytes []byte
	file := self.Open()
	readBytes, err := io.ReadAtLeast(file, headBytes, limitSize)
	if readBytes != limitSize {
		return headBytes, fmt.Errorf("error: failed to complete read: read ", readBytes, " out of ", limitSize, "bytes")
	} else {
		return headBytes, err
	}
}

// Need stream and zerocopy

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
	return self.Metadata().IsDir()
}

func (self Path) IsFile() bool {
	return self.Metadata().Mode().IsRegular()
}
