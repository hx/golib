package paths_test

import (
	"fmt"
	. "github.com/hx/golib/paths"
	. "github.com/hx/golib/testing"
	"runtime"
	"testing"
)

func pathTestFile() *Path {
	var _, thisFilePath, _, _ = runtime.Caller(0)
	return NewTree().Join(thisFilePath)
}

func TestPath_Join(t *testing.T) {
	root := NewTree()
	e := func(expected string, parts ...string) {
		if runtime.GOOS == "windows" {
			for i, v := range parts {
				parts[i] = reslash(v)
			}
			Equals(t, LocalSystem.Root()+reslash(expected[1:]), root.Join(parts...).String())
		} else {
			Equals(t, expected, root.Join(parts...).String())
		}
	}
	e("/foo", "foo")
	e("/foo/bar", "foo", "bar")
	e("/bar", "foo", "..", "bar")
	e("/baz", "foo", "/baz")
	e("/", "..")
	e("/", "..", "..")
}

func TestPath_Exists(t *testing.T) {
	Assert(t, pathTestFile().Exists(), "calling file should exist")
	Assert(t, !pathTestFile().Join("fqh394c8ol8i").Exists(), "gobbledygook file should not exist")
}

func TestPath_IsDir(t *testing.T) {
	Assert(t, pathTestFile().Parent().IsDir(), "calling file parent is a directory")
	Assert(t, !pathTestFile().IsDir(), "calling file is not a directory")
}

func TestPath_IsNonDir(t *testing.T) {
	Assert(t, pathTestFile().IsNonDir(), "calling file is not a directory")
	Assert(t, !pathTestFile().Parent().IsNonDir(), "calling file parent is a directory")
}

func TestPath_MustStat(t *testing.T) {
	stat := pathTestFile().MustStat()
	Equals(t, "path_test.go", stat.Name())

	notAFile := pathTestFile().Parent().Join("faliowxejea")
	err := func() (err interface{}) {
		defer func() { err = recover() }()
		notAFile.MustStat()
		return
	}()
	if runtime.GOOS == "windows" {
		Equals(t, fmt.Sprintf("CreateFile %s: The system cannot find the file specified.", notAFile), err.(error).Error())
	} else {
		Equals(t, fmt.Sprintf("lstat %s: no such file or directory", notAFile), err.(error).Error())
	}
}

func TestPath_ReadLink(t *testing.T) {
	if runtime.GOOS == "windows" {
		// The link in the repo isn't readable on Windows.
		// TODO: this test, on Windows
		return
	}
	target := pathTestFile()
	link := target.Parent().Join("path_test.link")
	Equals(t, target.String(), link.MustReadLink().String())
}
