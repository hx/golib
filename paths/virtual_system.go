package paths

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type VirtualSystem struct {
	rootDir  *virtualDir
	rootPath string
}

func NewVirtualSystem() *VirtualSystem {
	return &VirtualSystem{
		rootDir:  newVirtualDir(),
		rootPath: LocalSystem.Root(),
	}
}

func (v *VirtualSystem) Root() string { return v.rootPath }

func (v *VirtualSystem) Lstat(name string) (os.FileInfo, error) {
	entry := v.rootDir.resolve(name)
	if entry == nil {
		return nil, ErrPathNotFound
	}
	return entry, nil
}

func (v *VirtualSystem) Chmod(name string, mode os.FileMode) error {
	entry := v.rootDir.resolve(name)
	if entry == nil {
		return ErrPathNotFound
	}
	base := entry.entry()
	base.mode = mode.Perm() | base.mode.Type()
	return nil
}

func (v *VirtualSystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	entry := v.rootDir.resolve(name)
	if entry == nil {
		return ErrPathNotFound
	}
	base := entry.entry()
	base.accessed = atime
	base.modified = mtime
	return nil
}

func (v *VirtualSystem) CurrentUser() (*user.User, error) {
	return &user.User{
		Uid:      "0",
		Gid:      "0",
		Username: "root",
		Name:     "Root",
		HomeDir:  v.rootPath,
	}, nil
}

func (v *VirtualSystem) Getwd() (dir string, err error) { return v.rootPath, nil }

func (v *VirtualSystem) Glob(pattern string) (matches []string, err error) {
	panic("not implemented")
}

func (v *VirtualSystem) Join(elem ...string) string { return filepath.Join(elem...) }

func (v *VirtualSystem) MkdirAll(path string, perm os.FileMode) error {
	dir := v.rootDir
	parts := strings.Split(v.chompSeparator(path), string(os.PathSeparator))
	for i, part := range parts {
		if i == 0 && part == "" {
			continue
		}
		entry := dir.resolve(part)
		switch entry := entry.(type) {
		case nil:
			for _, part := range parts[i:] {
				child := newVirtualDir()
				child.name = part
				child.mode = child.mode.Type() | perm.Perm()
				child.parent = dir
				dir.children = append(dir.children, child)
				dir = child
			}
			return nil
		case *virtualSymlink:
			target := entry.resolveRecursive()
			switch target := target.(type) {
			case nil:
				return ErrBrokenLink
			case *virtualDir:
				dir = target
			default:
				return ErrNonDirectory
			}
		case *virtualDir:
			dir = entry
		default:
			return ErrNonDirectory
		}
	}
	return nil
}

func (v *VirtualSystem) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	entry := v.rootDir.resolve(name)
	if link, ok := entry.(*virtualSymlink); ok {
		entry = link.resolveRecursive()
	}
	if _, ok := entry.(*virtualDir); ok {
		return nil, ErrDirectory
	}
	file, isFile := entry.(*virtualFile)
	if !isFile && flag&(os.O_CREATE|os.O_APPEND|os.O_TRUNC) == 0 {
		return nil, ErrPathNotFound
	}
	if isFile && flag&os.O_EXCL != 0 {
		return nil, ErrFileExists
	}
	if isFile && flag&os.O_TRUNC != 0 {
		file.contents = []byte{}
	}
	if !isFile {
		dir, name, err := v.dirAndName(name)
		if err != nil {
			return nil, err
		}
		file = newVirtualFile(name, dir, perm)
		dir.children = append(dir.children, file)
	}
	f := &virtualOpenFile{file: file}
	if flag&os.O_APPEND != 0 {
		must1(f.Seek(0, 2))
	}
	return f, nil
}

func (v *VirtualSystem) ReadDir(name string) (entries []os.DirEntry, err error) {
	entry := v.rootDir.resolve(name)
	if link, ok := entry.(*virtualSymlink); ok {
		entry = link.resolveRecursive()
	}
	dir, isDir := entry.(*virtualDir)
	if !isDir {
		return nil, ErrNonDirectory
	}
	entries = make([]os.DirEntry, len(dir.children))
	for i, e := range dir.children {
		entries[i] = e
	}
	return
}

func (v *VirtualSystem) ReadFile(name string) (b []byte, err error) {
	file, err := v.OpenFile(name, 0, 0)
	if err != nil {
		return
	}
	return file.(*virtualOpenFile).file.contents, nil
}

func (v *VirtualSystem) Readlink(name string) (string, error) {
	entry := v.rootDir.resolve(name)
	if link, ok := entry.(*virtualSymlink); ok {
		return link.target, nil
	}
	return "", ErrNonLink
}

func (v *VirtualSystem) Remove(name string) error {
	entry := v.rootDir.resolve(name)
	if entry == nil {
		return ErrPathNotFound
	}
	entryBase := entry.entry()
	parent := entryBase.parent
	if parent == nil {
		return ErrNotWritable
	}
	for i, child := range parent.children {
		if child == entry {
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
			break
		}
	}
	entryBase.parent = nil
	return nil
}

func (v *VirtualSystem) RemoveAll(path string) error { return v.Remove(path) }

func (v *VirtualSystem) Rename(oldpath, newpath string) error {
	entry := v.rootDir.resolve(oldpath)
	if entry == nil {
		return ErrPathNotFound
	}
	entryBase := entry.entry()
	oldParent := entryBase.parent
	if oldParent == nil {
		return ErrNotWritable
	}
	newParent, newName, err := v.dirAndName(newpath)
	if err != nil {
		return err
	}
	_ = v.Remove(oldpath)
	_ = v.Remove(newpath)
	entryBase.name = newName
	entryBase.parent = newParent
	newParent.children = append(newParent.children, entry)
	return nil
}

func (v *VirtualSystem) SupportsSymlinks() bool { return LocalSystem.SupportsSymlinks() }

func (v *VirtualSystem) Symlink(oldname, newname string) error {
	dir, name, err := v.dirAndName(newname)
	if err != nil {
		return err
	}
	_ = v.Remove(newname)
	link := newVirtualSymlink(name, dir, 0644, oldname)
	dir.children = append(dir.children, link)
	return nil
}

func (v *VirtualSystem) dirAndName(path string) (dir *virtualDir, name string, err error) {
	splitAt := strings.LastIndexByte(path, os.PathSeparator)
	if splitAt == -1 {
		return nil, "", ErrInvalid
	}
	newParentEntry := v.rootDir.resolve(path[:splitAt])
	if link, ok := newParentEntry.(*virtualSymlink); ok {
		newParentEntry = link.resolveRecursive()
	}
	if newParentEntry == nil {
		return nil, "", ErrPathNotFound
	}
	newParent, isDir := newParentEntry.(*virtualDir)
	if !isDir {
		return nil, "", ErrNonDirectory
	}
	return newParent, path[splitAt+1:], nil
}

func (v *VirtualSystem) chompSeparator(path string) string {
	l := len(path)
	if path[l-1] == os.PathSeparator {
		return path[:l-1]
	}
	return path
}
