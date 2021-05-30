package trees

type nodeSet []*node

func (s nodeSet) height() (height int) {
	for _, n := range s {
		height += n.height()
	}
	return
}

func (s *nodeSet) add() (child *node) {
	child = new(node)
	*s = append(*s, child)
	return
}

func (s nodeSet) walk(level int, f walkFunc) (keepWalking bool) {
	keepWalking = true
	for _, n := range s {
		if !n.walk(level, f) {
			return false
		}
	}
	return true
}
