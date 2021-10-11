package paths

import (
	. "github.com/hx/golib/testing"
	"io"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"testing"
	"time"
)

func TestVirtualSystem_Root(t *testing.T) {
	root := NewVirtualSystem().Root()
	if runtime.GOOS == "windows" {
		Assert(t,
			regexp.MustCompile("^[A-Z]:\\\\$").MatchString(root),
			root+" should be a letter, colon, and backslash")
	} else {
		Equals(t, "/", root)
	}
}

func testSys() (sys *VirtualSystem, tree *Tree) {
	sys = NewVirtualSystem()
	tree = NewTreeWithSystem(sys)

	return
}

func TestVirtualSystem_Lstat(t *testing.T) {
	sys, tree := testSys()
	tree.Join("foo").MustWriteString("bar")

	stat, err := sys.Lstat(reslash("/foo"))
	Ok(t, err)
	Equals(t, "foo", stat.Name())
	Equals(t, false, stat.IsDir())
	Equals(t, sys.rootDir.children[0].entry().mode, stat.Mode())
	Equals(t, sys.rootDir.children[0].entry().modified, stat.ModTime())
	Equals(t, int64(3), stat.Size())
	Equals(t, nil, stat.Sys())
}

func TestVirtualSystem_Chmod(t *testing.T) {
	_, tree := testSys()
	foo := tree.Join("foo").MustTouch()

	Equals(t, fs.FileMode(0644), foo.MustStat().Mode())
	Ok(t, foo.Chmod(0123))
	Equals(t, fs.FileMode(0123), foo.MustStat().Mode())
}

func TestVirtualSystem_Chtimes(t *testing.T) {
	sys, tree := testSys()
	tree.Join("foo").MustTouch()

	newTime := time.Now().Add(time.Hour * -5)
	NotEquals(t, sys.rootDir.children[0].entry().accessed, newTime)
	NotEquals(t, sys.rootDir.children[0].entry().modified, newTime.Add(5))
	Ok(t, sys.Chtimes(reslash("/foo"), newTime, newTime.Add(5)))
	Equals(t, sys.rootDir.children[0].entry().accessed, newTime)
	Equals(t, sys.rootDir.children[0].entry().modified, newTime.Add(5))
}

func TestVirtualSystem_MkdirAll(t *testing.T) {
	sys, tree := testSys()
	Ok(t, sys.MkdirAll(reslash("/foo/bar/baz"), 0755))
	Assert(t, tree.Join(reslash("foo/bar/baz")).IsDir(), "dir should exist")
	Equals(t, 0755|fs.ModeDir, tree.Join(reslash("/foo/bar/baz")).MustStat().Mode())
	Ok(t, sys.MkdirAll(reslash("/foo/bar/baz"), 0755))
}

func TestVirtualSystem_OpenFile(t *testing.T) {
	sys, tree := testSys()

	t.Run("create", func(t *testing.T) {
		f, err := sys.OpenFile(reslash("/foo"), os.O_CREATE|os.O_EXCL, 0644)
		Ok(t, err)
		n, err := f.Write([]byte("hello"))
		Ok(t, err)
		Equals(t, 5, n)
		Ok(t, f.Close())
		Equals(t, "hello", tree.Join("foo").MustReadString())
	})

	t.Run("append", func(t *testing.T) {
		f, err := sys.OpenFile(reslash("/foo"), os.O_APPEND, 0644)
		Ok(t, err)
		n, err := f.Write([]byte(" world"))
		Ok(t, err)
		Equals(t, 6, n)
		Ok(t, f.Close())
		Equals(t, "hello world", tree.Join("foo").MustReadString())
	})

	t.Run("truncate", func(t *testing.T) {
		f, err := sys.OpenFile(reslash("/foo"), os.O_TRUNC, 0644)
		Ok(t, err)
		n, err := f.Write([]byte("bye"))
		Ok(t, err)
		Equals(t, 3, n)
		Ok(t, f.Close())
		Equals(t, "bye", tree.Join("foo").MustReadString())
	})

	t.Run("read", func(t *testing.T) {
		f, err := sys.OpenFile(reslash("/foo"), os.O_RDONLY, 0644)
		Ok(t, err)
		b, err := io.ReadAll(f)
		Ok(t, err)
		Ok(t, f.Close())
		Equals(t, "bye", string(b))
	})
}

func TestVirtualSystem_ReadDir(t *testing.T) {
	sys, tree := testSys()

	tree.Join("foo").MustTouch()
	tree.Join("bar").MustMake()

	entries, err := sys.ReadDir(reslash("/"))
	Ok(t, err)
	Equals(t, 2, len(entries))

	foo := entries[0].(*virtualFile)
	Equals(t, "foo", foo.name)

	bar := entries[1].(*virtualDir)
	Equals(t, "bar", bar.name)
}

func TestVirtualSystem_ReadFile(t *testing.T) {
	sys, tree := testSys()
	tree.Join("foo").MustWriteString("bar")
	b, err := sys.ReadFile("foo")
	Ok(t, err)
	Equals(t, "bar", string(b))
}

func TestVirtualSystem_Remove(t *testing.T) {
	sys, tree := testSys()
	tree.Join("foo").MustTouch()
	tree.Join("bar").MustMake()
	Equals(t, 2, len(sys.rootDir.children))
	Ok(t, sys.Remove("foo"))
	Equals(t, 1, len(sys.rootDir.children))
	Ok(t, sys.Remove("bar"))
	Equals(t, 0, len(sys.rootDir.children))
}

func TestVirtualSystem_Rename(t *testing.T) {
	sys, tree := testSys()
	tree.Join("foo").MustMake()
	Ok(t, sys.Rename(reslash("/foo"), reslash("/bar")))
	Equals(t, 1, len(sys.rootDir.children))
	Equals(t, "bar", sys.rootDir.children[0].entry().name)
}

func TestVirtualSystem_Readlink(t *testing.T) {
	sys, tree := testSys()
	tree.Join("foo").MustMake()
	Ok(t, sys.Symlink(reslash("/foo"), reslash("/foo/back2foo")))
	link, err := sys.Readlink(reslash("/foo/back2foo"))
	Ok(t, err)
	Equals(t, reslash("/foo"), link)
}

func TestVirtualSystem_Symlink(t *testing.T) {
	sys, tree := testSys()
	tree.Join("foo").MustMake()
	Ok(t, sys.Symlink(reslash("/foo"), reslash("/foo/back2foo")))
	Equals(t, reslash("/foo"), sys.rootDir.children[0].(*virtualDir).children[0].(*virtualSymlink).target)
}
