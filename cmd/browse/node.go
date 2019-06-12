package browse

type node struct {
	name     string
	isDir    bool
	size     int64
	children []node
	parent   *node
}

func (n *node) has(name string) bool {
	for _, c := range n.children {
		if c.name == name {
			return true
		}
	}
	return false
}

func (n *node) getChild(name string) (*node, bool) {
	for i, c := range n.children {
		if c.name == name {
			// Get a reference to array element directly, since range returns copies.
			return &n.children[i], true
		}
	}
	return nil, false
}
