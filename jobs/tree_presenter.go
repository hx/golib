package jobs

import "github.com/hx/golib/trees"

// TreePresenter wraps a Job, and uses a Formatter to display each of Job's descendants in a node of a trees.Node.
type TreePresenter struct {
	formatter Formatter
	tree      *TreeBinding
}

// NewTreePresenter creates TreePresenter.
func NewTreePresenter(job Job, node trees.Node, formatter Formatter) *TreePresenter {
	return &TreePresenter{
		tree:      NewTreeBinding(job, node),
		formatter: formatter,
	}
}

// Run implements Job.
func (t *TreePresenter) Run(ctx *Context) (err error) {
	for event := range ctx.Run(t.tree) {
		if node := t.tree.Node(event.Job()); node != nil {
			node.Update(t.formatter(event))
		}
		if event, ok := event.(*EventFinished); ok {
			err = event.Error()
		}
	}
	return
}
