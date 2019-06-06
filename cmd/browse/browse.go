package browse

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/robinmitra/forest/formatter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type options struct {
	tree bool
	root string
}

func (o *options) initialise(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		r, _ := regexp.Compile("/$")
		o.root = r.ReplaceAllString(args[0], "")
	} else {
		o.root = "."
	}
	if tree, _ := cmd.Flags().GetBool("tree"); tree {
		o.tree = tree
	}
}

func (o *options) validate() {
	if err := o.validatePath(os.Stat(o.root)); err != nil {
		log.Fatal(err)
	}
}

func (o *options) validatePath(info os.FileInfo, err error) error {
	if os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Directory \"%s\" does not exist", o.root))
	}
	return err
}

func (o *options) run() {
	if !o.tree {
		log.Fatal("Unknown display mode")
		return
	}
	node := buildFileTree(o.root)
	renderTree(node)
}

var cmd = &cobra.Command{
	Use:   "browse",
	Short: "Interactively browse directories and files",
}

func NewInteractiveCmd() *cobra.Command {
	o := options{}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		o.initialise(cmd, args)
		o.validate()
		o.run()
	}

	cmd.Flags().BoolVarP(
		&o.tree,
		"tree",
		"t",
		true,
		"browse the file tree",
	)

	return cmd
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

	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}

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
