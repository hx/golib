package paths

import (
	"io"
	"os"
	"os/user"
	"time"
)

type File interface {
	io.ReadCloser
	io.WriteSeeker
}

type System interface {
	Chmod(name string, mode os.FileMode) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
	CurrentUser() (*user.User, error)
	Getwd() (dir string, err error)
	Glob(pattern string) (matches []string, err error)
	Join(elem ...string) string
	Lstat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	ReadDir(name string) ([]os.DirEntry, error)
	ReadFile(name string) ([]byte, error)
	Readlink(name string) (string, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldpath, newpath string) error
	Root() string
	SupportsSymlinks() bool
	Symlink(oldname, newname string) error
}
