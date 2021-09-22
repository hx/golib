package paths

import "sync/atomic"

type Tree struct {
	*Path
	sys       System
	listeners map[uint64]Listener
}

func NewTree() (t *Tree) { return NewTreeWithSystem(local{}) }

func NewTreeWithSystem(system System) (t *Tree) {
	t = &Tree{
		sys:       system,
		listeners: make(map[uint64]Listener),
	}
	t.Path = &Path{t.sys.Root(), t}
	return

}

var nextID = new(uint64)

func (t *Tree) Subscribe(listener Listener) (unsubscribe func()) {
	id := atomic.AddUint64(nextID, 1)
	t.listeners[id] = listener
	return func() { delete(t.listeners, id) }
}

func (t *Tree) dispatch(event Event) {
	for _, l := range t.listeners {
		l(event)
	}
}
