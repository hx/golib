package trees

// Node is a single node in a Tree (including the Tree itself).
type Node interface {
	// AddChild adds a new child Node to the receiver.
	AddChild() (child Node)

	// Update displays content on the Node's line, replacing any content set by previous calls to Update. There must
	Update(line string)
}

type node struct {
	children      nodeSet
	content       string
	contentHeight int
	tree          *Tree
}

// AddChild implements Node.
func (n *node) AddChild() (child Node) {
	n.tree.mutex.Lock()
	ch := n.children.add()
	ch.tree = n.tree
	n.tree.mutex.Unlock()
	return ch
}

// Update implements Node.
func (n *node) Update(content string) {
	n.tree.mutex.Lock()
	n.tree.render(n, content)
	n.tree.mutex.Unlock()
}

// height returns the total height of nodes in the receiver, including itself and all its descendents.
func (n *node) height() (height int) {
	return n.contentHeight + n.children.height()
}

type walkFunc func(node *node, level int) bool

func (n *node) walk(level int, f walkFunc) (keepWalking bool) {
	if f(n, level) {
		return n.children.walk(level+1, f)
	}
	return false
}
