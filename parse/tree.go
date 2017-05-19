package parse

import (
	"fmt"
	"io"
	"strings"
)

// Kind enumerates what kind of information a Node holds.
type Kind int

const (
	KindNone       Kind = iota // Any node other than those listed below.
	KindType                   // A Node representing a type definition.
	KindUnionType              // A Node representing a union.
	KindMember                 // A Node representing a variable of a specific type (TypeName) defined elsewhere.
	KindBaseMember             // A Node representing a variable of a Go builtin type.
	KindModule                 // A Node representing a module.
	KindConst                  // A Node representing a const(ant).
	KindEnum                   // A Node representing an enum.
)

// Node represents an item in an IDL specification and form trees with one
// node at the top. Nodes have a pointer to a parent Node (or nil for the
// topmost Node) and zero or more pointers to child Nodes.
//
// The Name of a child is guaranteed to be unique among its siblings.
// The order of the children is preserved. The first child added is the
// first in the collection returned by Children().
type Node struct {
	Name   string // The name of the node. May be an empty string.
	Parent *Node  // The parent Node. May be nil.

	Kind     Kind
	TypeName string // If Kind == KindMember or KindConst, this is the type of the variable.
	Length   string // "": not a sequence
	Value    string // If Kind == KindConst, this is the value.

	children []*Node
}

// Child returns a pointer to the name named "name" or nil (if no such child
// exist.
func (n *Node) Child(name string) *Node {
	for _, c := range n.children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// Add ensures that "n" has a children named like the children in the argument
// and that they in turn have the same children.
//
// Put differently, Add adds the given children while ensuring that no two
// children share the same name.
func (n *Node) Add(c ...*Node) {
	for _, child := range c {
		if child == nil {
			return
		}

		if existingChild := n.Child(child.Name); existingChild != nil {
			// A child with that name already exist. Move all children of "c" onto
			// "existingChild".
			existingChild.Add(child.Children()...)
			// Empty the child we're discarding of parents and children
			child.Parent = nil
			child.children = nil
		} else {
			// There is no child with this name. Add the child.
			n.children = append(n.children, child)
			child.Parent = n
		}
	}
}

// Children returns the array of Children of this Node.
func (n *Node) Children() []*Node {
	return n.children
}

// Traverse will apply "test()" to each node in thee starting at "node" and,
// for each node where test() returns true, apply do().
func (n *Node) Traverse(test func(n *Node) bool, do func(*Node)) {
	if test(n) {
		do(n)
	}

	// Ensure that traverse maintains the order in which the children were added.
	for index := range n.children {
		n.children[index].Traverse(test, do)
	}
}

// Dumps a textual representation of the tree to the io.Writer provided,
// for example "os.Stdout".
func (n *Node) Dump(out io.Writer) {
	// Define the function so that the variable can be used within the function
	// itself.
	var f func(n *Node, indent int, out io.Writer)

	f = func(n *Node, indent int, out io.Writer) {
		if indent > 100 {
			return // Emergency bailout.
		}
		fmt.Fprintf(out, "%s\"%s\" %p %+v\n", strings.Repeat(" ", indent), n.Name, n, *n)
		for _, c := range n.children {
			f(c, indent+1, out)
		}
	}
	f(n, 0, out)
}

// FullPathName returns the Names of the Nodes starting at "n" and ending at
// the top node (the one with no parent) or as soon as a node with an empty
// name is found in the chain towards the top node, and inserts the "divider"
// between each.
func (n *Node) FullPathName(divider string) string {
	if len(n.Name) == 0 {
		return ""
	}

	s := ""
	if n.Parent != nil && len(n.Parent.Name) != 0 {
		s += n.Parent.FullPathName(divider) + divider
	}

	return s + n.Name
}

// Find returns the Node with the "name" at or below the node.
func (n *Node) Find(name string) *Node {
	var node *Node
	n.Traverse(func(n *Node) bool {
		return n.Name == name
	}, func(n *Node) {
		node = n
	})
	return node
}

// Equal compares two trees and returns true if they contain nodes with the
// same names, typename, children, etc. They need not point to the same Nodes
// as long as the contents of the Nodes are the same.
func (n *Node) Equal(node *Node) bool {
	if node == nil {
		return false
	}
	if n.Name != node.Name {
		return false
	}
	if n.TypeName != node.TypeName {
		return false
	}
	if n.Length != node.Length {
		return false
	}
	if n.Kind != node.Kind {
		return false
	}
	if n.Value != node.Value {
		return false
	}
	if (n.Parent == nil && node.Parent != nil) || (n.Parent != nil && node.Parent == nil) {
		return false
	}
	if len(n.children) != len(node.children) {
		return false
	}
	for index := range n.children {
		if !n.children[index].Equal(node.children[index]) {
			return false
		}
	}
	return true
}

// getTop returns the topmost Node in the tree "n" sits in.
func (n *Node) GetTop() *Node {
	top := n
	for top.Parent != nil {
		top = top.Parent
	}
	return top
}