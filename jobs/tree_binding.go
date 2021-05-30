package jobs

import (
	"github.com/hx/golib/trees"
	"sync"
)

// TreeBinding creates a trees.Node on the given trees.Node for every descendent Job spawned by the given Job,
// and provides them through the Node function. TreeBinding implements Job, and can be run like any other Job.
// A zero TreeBinding is not valid; use NewTreeBinding.
type TreeBinding struct {
	job  Job
	lock sync.RWMutex
	// A sync.Map doesn't quite work here, because of the delay between storage and availability
	nodes map[Job]trees.Node
}

// NewTreeBinding creates a new TreeBinding for the given job and node.
func NewTreeBinding(job Job, node trees.Node) *TreeBinding {
	return &TreeBinding{
		job:   job,
		nodes: map[Job]trees.Node{job: node},
	}
}

// Run implements Job.
func (t *TreeBinding) Run(ctx *Context) (err error) {
	for event := range ctx.Run(t.job) {
		switch event := event.(type) {
		case *EventQueued:
			if parent := t.Node(event.Parent()); parent != nil {
				t.lock.Lock()
				t.nodes[event.Job()] = parent.AddChild()
				t.lock.Unlock()
			}
		case *EventFinished:
			if event.Error() != nil {
				err = event.Error()
			}
		}
	}
	return
}

// Node returns the trees.Node for the given job, if the job has already been queued as the TreeBinding's job or one
// of its descendents. Otherwise, Node returns nil.
func (t *TreeBinding) Node(job Job) (node trees.Node) {
	t.lock.RLock()
	node = t.nodes[job]
	t.lock.RUnlock()
	return
}
