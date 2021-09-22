package paths

import (
	"os"
	"path/filepath"
	"time"
)

type Path struct {
	path string
	tree *Tree
}

func (p *Path) Join(a ...string) *Path {
	n := *p

	for _, str := range a {
		if filepath.IsAbs(str) {
			n.path = filepath.Clean(str)
		} else {
			n.path = filepath.Join(n.path, str)
		}
	}

	return &n
}

func (p *Path) Parent() *Path                { return p.Join("..") }
func (p *Path) Stat() (os.FileInfo, error)   { return p.tree.sys.Lstat(p.path) }
func (p *Path) MustStat() os.FileInfo        { return must1(p.Stat()).(os.FileInfo) }
func (p *Path) String() string               { return p.path }
func (p *Path) Base() string                 { return filepath.Base(p.path) }
func (p *Path) Extension() string            { return filepath.Ext(p.path) }
func (p *Path) Chmod(mode os.FileMode) error { return p.tree.sys.Chmod(p.path, mode) }
func (p *Path) MustChmod(mode os.FileMode)   { must(p.Chmod(mode)) }

func (p *Path) Exists() bool {
	_, err := p.Stat()
	return err == nil
}

func (p *Path) Touch() error {
	if p.Exists() {
		now := time.Now()
		return p.tree.sys.Chtimes(p.path, now, now)
	}

	file, err := p.Create()
	if err == nil {
		err = file.Close()
	}
	return err
}
func (p *Path) MustTouch() *Path {
	must(p.Touch())
	return p
}

func (p *Path) IsDir() bool {
	s, err := p.Stat()
	return err == nil && s.IsDir()
}

func (p *Path) IsNonDir() bool {
	s, err := p.Stat()
	return err == nil && !s.IsDir()
}

func (p *Path) Wd() (*Path, error) {
	s, err := p.tree.sys.Getwd()
	if err != nil {
		return nil, err
	}
	return p.Join(s), nil
}
func (p *Path) MustWd() *Path { return must1(p.Wd()).(*Path) }

func (p *Path) UserHome() (*Path, error) {
	u, err := p.tree.sys.CurrentUser()
	if err != nil {
		return nil, err
	}
	return p.Join(u.HomeDir), nil
}

func (p *Path) Size() (int64, error) {
	stat, err := p.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}
func (p *Path) MustSize() int64 { return must1(p.Size()).(int64) }

func (p *Path) SizeIfExists() (int64, error) {
	if !p.Exists() {
		return -1, nil
	}
	return p.Size()
}
func (p *Path) MustSizeIfExists() int64 { return must1(p.SizeIfExists()).(int64) }

func (p *Path) ReadLink() (*Path, error) {
	f, err := p.tree.sys.Readlink(p.path)
	if err != nil {
		return nil, err
	}
	return p.Parent().Join(f), nil
}
func (p *Path) MustReadLink() *Path { return must1(p.ReadLink()).(*Path) }

func (p *Path) IsEmpty() (isEmpty bool, err error) {
	if p.IsDir() {
		var all []os.DirEntry
		all, err = p.tree.sys.ReadDir(p.path)
		isEmpty = len(all) == 0
		return
	}
	if p.IsNonDir() {
		var info os.FileInfo
		info, err = p.tree.sys.Lstat(p.path)
		if info != nil {
			isEmpty = info.Size() > 0
		}
		return
	}
	return true, nil
}
func (p *Path) MustIsEmpty() bool { return must1(p.IsEmpty()).(bool) }
