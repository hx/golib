package paths

import (
	"bytes"
	"io"
	"os"
)

func (p *Path) CreateMode(mode os.FileMode) (file File, err error) {
	if err = p.Parent().Make(); err != nil {
		return
	}
	return p.tree.sys.OpenFile(p.path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, mode)
}

func (p *Path) Create() (File, error)                { return p.CreateMode(0644) }
func (p *Path) MustCreate() File                     { return must1(p.Create()).(File) }
func (p *Path) MustCreateMode(mode os.FileMode) File { return must1(p.CreateMode(mode)).(File) }

func (p *Path) AppendMode(mode os.FileMode) (file File, err error) {
	if err = p.Parent().Make(); err != nil {
		return
	}
	return p.tree.sys.OpenFile(p.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, mode)
}

func (p *Path) Append() (File, error)                { return p.AppendMode(0644) }
func (p *Path) MustAppend() File                     { return must1(p.Append()).(File) }
func (p *Path) MustAppendMode(mode os.FileMode) File { return must1(p.AppendMode(mode)).(File) }

func (p *Path) Open() (File, error) { return p.tree.sys.OpenFile(p.path, os.O_RDONLY, 0) }
func (p *Path) MustOpen() File      { return must1(p.Open()).(File) }

func (p *Path) WriteFrom(reader io.Reader) (err error) {
	writer := p.WriteCloser()
	_, err = io.Copy(writer, reader)
	if err == nil {
		err = writer.Close()
	} else {
		_ = writer.Close()
	}
	return
}
func (p *Path) MustWriteFrom(reader io.Reader) { must(p.WriteFrom(reader)) }
func (p *Path) WriteBytes(b []byte) error      { return p.WriteFrom(bytes.NewReader(b)) }
func (p *Path) MustWriteBytes(b []byte)        { must(p.WriteBytes(b)) }
func (p *Path) WriteString(s string) error     { return p.WriteBytes([]byte(s)) }
func (p *Path) MustWriteString(s string)       { must(p.WriteString(s)) }

func (p *Path) WriteBytesUnlessEqual(b []byte) error {
	if !p.Exists() {
		return p.WriteBytes(b)
	}
	eq, err := p.BytesAreEqualToBytes(b)
	if err == nil && !eq {
		err = p.WriteBytes(b)
	}
	return err
}
func (p *Path) MustWriteBytesUnlessEqual(b []byte)    { must(p.WriteBytesUnlessEqual(b)) }
func (p *Path) WriteStringUnlessEqual(s string) error { return p.WriteBytesUnlessEqual([]byte(s)) }
func (p *Path) MustWriteStringUnlessEqual(s string)   { must(p.WriteStringUnlessEqual(s)) }

func (p *Path) ReadBytesIfExists() (b []byte, err error) {
	if !p.Exists() {
		return
	}
	return p.ReadBytes()
}
func (p *Path) MustReadBytesIfExists() []byte    { return must1(p.ReadBytesIfExists()).([]byte) }
func (p *Path) ReadBytes() (b []byte, err error) { return p.tree.sys.ReadFile(p.path) }
func (p *Path) MustReadBytes() []byte            { return must1(p.ReadBytes()).([]byte) }

func (p *Path) ReadString() (string, error) {
	b, err := p.ReadBytes()
	return string(b), err
}
func (p *Path) ReadStringIfExists() (str string, err error) {
	if !p.Exists() {
		return
	}
	return p.ReadString()
}
func (p *Path) MustReadString() string         { return must1(p.ReadString()).(string) }
func (p *Path) MustReadStringIfExists() string { return must1(p.ReadStringIfExists()).(string) }

func (p *Path) ReadTo(writer io.Writer) (err error) {
	file, err := p.Open()
	if err != nil {
		return
	}
	_, err = io.Copy(writer, file)
	if err == nil {
		err = file.Close()
	} else {
		_ = file.Close()
	}
	return
}
func (p *Path) MustReadTo(writer io.Writer) { must(p.ReadTo(writer)) }
func (p *Path) ReadToIfExists(writer io.Writer) error {
	if !p.Exists() {
		return nil
	}
	return p.ReadTo(writer)
}
func (p *Path) MustReadToIfExists(writer io.Writer) { must(p.ReadToIfExists(writer)) }
