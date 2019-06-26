package browse

import (
	"testing"
)

func TestCanAddAndRetrieveChild(t *testing.T) {
	n := node{}
	c1 := node{name: "C1", size: 100, isDir: false}
	c2 := node{name: "C2", size: 200, isDir: false}
	c3 := node{name: "C3", isDir: true}

	n.addChild(&c1)
	n.addChild(&c2)
	n.addChild(&c3)

	if len(n.children) != 3 {
		t.Fatalf("Expected node to have 2 children, found %d", len(n.children))
	}
	if n.size != 300 {
		t.Fatalf("Expected node to have size of %d, found %d", 300, n.size)
	}
	if !n.hasChild("C1") || !n.hasChild("C2") {
		t.Fatalf("Expected node to have children C1 and C2")
	}
	if n.hasChild("C4") {
		t.Fatalf("Expected node to not have child C4")
	}
	if c, ok := n.getChild("C1"); &c1 != c || !ok {
		t.Fatalf("Expected node to have child C1")
	}
	if c, ok := n.getChild("C4"); c != nil || ok {
		t.Fatalf("Expected node to not have child C4")
	}
}

func TestCalculateSizeCorrectly(t *testing.T) {
	n := node{name: "R"}
	c1 := node{name: "C1", size: 100, isDir: false}
	c2 := node{name: "C2", size: 200, isDir: false}
	c3 := node{name: "C3", isDir: true}
	c31 := node{name: "C31", size: 300, isDir: false}
	c32 := node{name: "C32", size: 400, isDir: false}
	c33 := node{name: "C33", isDir: true}
	c331 := node{name: "C331", size: 500, isDir: false}
	c332 := node{name: "C332", size: 600, isDir: false}
	c333 := node{name: "C333", isDir: true}

	n.addChild(&c1)
	n.addChild(&c2)
	n.addChild(&c3)
	c3.addChild(&c31)
	c3.addChild(&c32)
	c3.addChild(&c33)
	c33.addChild(&c331)
	c33.addChild(&c332)
	c33.addChild(&c333)

	c33.recalculateSize()
	c3.recalculateSize()
	n.recalculateSize()

	if c33.size != 1100 {
		t.Fatalf("Expected child node c33 to have size of %d, found %d", 1100, c33.size)
	}
	if c3.size != 1800 {
		t.Fatalf("Expected child node c3 to have size of %d, found %d", 1800, c33.size)
	}
	if n.size != 2100 {
		t.Fatalf("Expected root node to have size of %d, found %d", 2100, c33.size)
	}
}
