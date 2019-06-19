package browse

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/robinmitra/forest/formatter"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func AggregateChildrenSize(n *node) int64 {
	var s int64
	for _, c := range n.children {
		s += c.size
	}
	return s
}

func buildNodesFromPath(n *node, path string, info os.FileInfo) {
	nodeNames := strings.Split(path, "/")
	currNodeName := nodeNames[0]
	nestedNodeNames := nodeNames[1:]
	// Last or trailing node
	if len(nestedNodeNames) == 0 {
		newNode := node{name: currNodeName, parent: n}
		if info.IsDir() {
			newNode.isDir = true
		} else {
			newNode.isDir = false
			newNode.size = info.Size()
			n.size += info.Size()
		}
		n.children = append(n.children, newNode)
	} else {
		if existingNode, ok := n.getChild(currNodeName); ok {
			buildNodesFromPath(existingNode, strings.Join(nestedNodeNames, "/"), info)
		} else {
			newNode := node{name: currNodeName, isDir: true, parent: n}
			buildNodesFromPath(&newNode, strings.Join(nestedNodeNames, "/"), info)
			n.children = append(n.children, newNode)
		}
		n.size = AggregateChildrenSize(n)
	}
}

func processFile(node *node, rootPath string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if path == "." || path == rootPath {
			node.name = name
			return nil
		}
		if rootPath != "." {
			buildNodesFromPath(node, strings.Replace(path, rootPath+"/", "", 1), info)
		} else {
			buildNodesFromPath(node, path, info)
		}
		return nil
	}
}

func buildFileTree(root string) *node {
	rootName := root
	if root != "." {
		path := strings.Split(root, "/")
		rootName = path[len(path)-1]
	}
	rootNode := node{name: rootName, isDir: true}
	if err := filepath.Walk(root, processFile(&rootNode, root)); err != nil {
		log.Fatal(err)
	}
	return &rootNode
}

func renderTree(n *node) {
	getNodeText := func(n *node) string {
		return fmt.Sprintf("%s (%s, %d)", n.name, formatter.HumaniseStorage(n.size), len(n.children))
	}

	root := tview.NewTreeNode(getNodeText(n)).SetReference(n).SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	addChildren := func(n *tview.TreeNode) {
		refNode := n.GetReference().(*node)
		if len(refNode.children) > 0 {
			for i, c := range refNode.children {
				cNode := &refNode.children[i]
				childNode := tview.NewTreeNode(getNodeText(cNode)).SetReference(cNode)
				if c.isDir {
					childNode.SetColor(tcell.ColorGreen)
				}
				n.AddChild(childNode)
			}
		}
	}

	addChildren(root)

	tree.SetSelectedFunc(func(n *tview.TreeNode) {
		refNode := n.GetReference()
		if refNode == nil {
			// Selecting the root node does nothing.
			return
		}
		children := n.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			addChildren(n)
		} else {
			// Collapse if visible, expand if collapsed.
			n.SetExpanded(!n.IsExpanded())
		}
	})

	app := tview.NewApplication().SetRoot(tree, true)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'o' {
			return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
		}
		return event
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func debugTree(n *node, spacer string) {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s", spacer, n.name)
	fmt.Fprintf(&b, " (")
	fmt.Fprintf(&b, "dir: %t", n.isDir)
	if n.parent != nil {
		fmt.Fprintf(&b, ", parent: %s", n.parent.name)
	}
	fmt.Fprintf(&b, ", size: %d", n.size)
	fmt.Fprintf(&b, ")")
	fmt.Println(b.String())

	if n.isDir && len(n.children) > 0 {
		for i, _ := range n.children {
			debugTree(&n.children[i], fmt.Sprintf("%s-", spacer))
		}
	}
}
