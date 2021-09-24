package paths

import (
	"os"
	"os/user"
	"path/filepath"
	"time"
)

type local struct{}

var LocalSystem System = local{}

func (l local) Chmod(name string, mode os.FileMode) error         { return os.Chmod(name, mode) }
func (l local) CurrentUser() (*user.User, error)                  { return user.Current() }
func (l local) Getwd() (dir string, err error)                    { return os.Getwd() }
func (l local) Glob(pattern string) (matches []string, err error) { return filepath.Glob(pattern) }
func (l local) Join(elem ...string) string                        { return filepath.Join(elem...) }
func (l local) Lstat(name string) (os.FileInfo, error)            { return os.Lstat(name) }
func (l local) MkdirAll(path string, perm os.FileMode) error      { return os.MkdirAll(path, perm) }
func (l local) ReadDir(name string) ([]os.DirEntry, error)        { return os.ReadDir(name) }
func (l local) ReadFile(name string) ([]byte, error)              { return os.ReadFile(name) }
func (l local) Readlink(name string) (string, error)              { return os.Readlink(name) }
func (l local) Remove(name string) error                          { return os.Remove(name) }
func (l local) RemoveAll(path string) error                       { return os.RemoveAll(path) }
func (l local) Rename(oldpath, newpath string) error              { return os.Rename(oldpath, newpath) }
func (l local) Symlink(oldname, newname string) error             { return os.Symlink(oldname, newname) }

func (l local) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (l local) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}
