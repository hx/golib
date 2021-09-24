package paths

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

type virtualFile struct {
	*virtualEntryBase
	contents []byte
}

func newVirtualFile(name string, parent *virtualDir, perm os.FileMode) *virtualFile {
	now := time.Now()
	return &virtualFile{
		virtualEntryBase: &virtualEntryBase{
			name:     name,
			mode:     perm.Perm(),
			accessed: now,
			modified: now,
			parent:   parent,
		},
	}
}

func (f *virtualFile) Size() int64                { return int64(len(f.contents)) }
func (f *virtualFile) IsDir() bool                { return false }
func (f *virtualFile) Info() (fs.FileInfo, error) { return f, nil }

type virtualOpenFile struct {
	file   *virtualFile
	offset int
}

func (f *virtualOpenFile) Read(p []byte) (n int, err error) {
	if f.offset == len(f.file.contents) {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return
	}
	n = copy(p, f.file.contents[f.offset:])
	f.offset += n
	f.file.accessed = time.Now()
	return
}

func (f *virtualOpenFile) Close() error { return nil }

func (f *virtualOpenFile) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	newSize := f.offset + len(p)
	if newSize > len(f.file.contents) {
		f.file.contents = append(f.file.contents[:f.offset], p...)
	} else {
		copy(f.file.contents[f.offset:], p)
	}
	n = len(p)
	f.offset += n
	now := time.Now()
	f.file.accessed = now
	f.file.modified = now
	return
}

func (f *virtualOpenFile) Seek(offset int64, whence int) (int64, error) {
	length := len(f.file.contents)
	switch whence {
	case 0:
		f.offset = int(offset)
	case 1:
		f.offset += int(offset)
	case 2:
		f.offset = length - int(offset)
	default:
		return 0, fmt.Errorf("%d is not valid for argument whence", whence)
	}
	if f.offset < 0 {
		f.offset = 0
	}
	if f.offset > length {
		f.offset = length
	}
	return int64(f.offset), nil
}
