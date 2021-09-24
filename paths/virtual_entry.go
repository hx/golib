package paths

import (
	"io/fs"
	"time"
)

type virtualEntry interface {
	fs.FileInfo
	fs.DirEntry
	entry() *virtualEntryBase
}

type virtualEntryBase struct {
	name     string
	mode     fs.FileMode
	accessed time.Time
	modified time.Time
	parent   *virtualDir
}

func (e *virtualEntryBase) entry() *virtualEntryBase { return e }
func (e *virtualEntryBase) Name() string             { return e.name }
func (e *virtualEntryBase) Type() fs.FileMode        { return e.mode }
func (e *virtualEntryBase) Mode() fs.FileMode        { return e.mode }
func (e *virtualEntryBase) ModTime() time.Time       { return e.modified }
func (e *virtualEntryBase) Sys() interface{}         { return nil }
