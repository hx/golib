package paths

import "io"

type RewriteFileEvent struct{ FileEvent }
type CreateFileEvent struct{ FileEvent }

type lazyWriteCloser struct {
	path    *Path
	existed bool
	written int
	file    io.WriteCloser
}

func (p *Path) WriteCloser() io.WriteCloser { return &lazyWriteCloser{path: p} }

func (l *lazyWriteCloser) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	if l.file == nil {
		l.existed = l.path.IsNonDir()
		if l.file, err = l.path.Create(); err != nil {
			return
		}
	}
	n, err = l.file.Write(p)
	l.written += n
	return
}

func (l *lazyWriteCloser) Close() error {
	file := l.file
	if file == nil {
		return nil
	}
	l.file = nil
	if l.existed {
		l.path.tree.dispatch(RewriteFileEvent{newFileEvent(l.path, int64(l.written))})
	} else {
		l.path.tree.dispatch(CreateFileEvent{newFileEvent(l.path, int64(l.written))})
	}
	return file.Close()
}
