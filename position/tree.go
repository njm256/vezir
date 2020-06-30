package position

type gameTree struct {
	metaData
	root *treeNode //should this be inherited rather than composed?
}

type metaData struct {
	wPlayer string
	bPlayer string
}

type treeNode struct {
	Game
	parent   *treeNode
	children []*treeNode
}

func (node treeNode) reps(s State) int {
	if node.prev == nil {
		count := 0
		n := &node
		for n.prev == nil {
			if n.state == s {
				count++
			}
			n = n.parent
		}
		return count + n.reps(s)
	}

	return node.prev[s.hashCode()]
}
