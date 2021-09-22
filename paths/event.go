package paths

type Event interface{ Path() *Path }

type event struct{ path *Path }

func newEvent(path *Path) Event { return &event{path} }

func (e *event) Path() *Path { return e.path }

type FileEvent interface {
	Event
	Size() int64
}

type fileEvent struct {
	event
	size int64
}

func newFileEvent(path *Path, size int64) FileEvent { return &fileEvent{event{path}, size} }

func (f *fileEvent) Size() int64 { return f.size }

type TargetEvent interface {
	Event
	Target() *Path
}

type targetEvent struct {
	event
	target *Path
}

func newTargetEvent(path *Path, target *Path) TargetEvent { return &targetEvent{event{path}, target} }

func (t *targetEvent) Target() *Path { return t.target }
