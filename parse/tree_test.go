package parse

import (
	"testing"
	"bytes"
	"strings"
)

func TestFreshNode(t *testing.T) {
	n := Node{}
	if len(n.Children()) != 0  {
		t.Error("Did not expect any children.")
	}
	if n.Parent != nil {
		t.Error("Did not expect a parent.")
	}
}

func TestNilChild(t *testing.T) {
	a := &Node{}
	a.Add(nil)
	if len(a.Children()) != 0 {
		t.Error("Expected no children.")
	}
}

func TestSmallTree(t *testing.T) {
	// a - b
	a := &Node{Name: "A"}
	b := &Node{Name: "B"}
	a.Add(b)

	if len(a.Children()) != 1 {
		t.Error("Expected one child.")
	}
	child := a.Child(b.Name)
	if child == nil {
		t.Error("Expected \"b\" to be a child of \"a\".")
	} else if child.Parent != a {
		t.Error("Expected the parent of \"b\" to be \"a\", but it was ", child.Parent)
	}
}

func TestChildrenWithSameName(t *testing.T) {
	a := &Node{Name: "A"}
	b := &Node{Name: "X"}
	c := &Node{Name: "X"}
	a.Add(b)
	a.Add(c)

	if len(a.Children()) != 1 {
		t.Error("Expected one child (b).")
	}
	if c.Parent != nil {
		t.Error("Expected the parent of \"c\" to remain nil.")
	}
}

func TestChildrenWithSameNameInTwoLevels(t *testing.T) {
	a := &Node{}
	b := &Node{}
	c := &Node{}
	d := &Node{}
	e := &Node{}
	f := &Node{Name: "f"}
	g := &Node{Name: "g"}

	// a - b - c - f
	//     d - e - g
	a.Add(b)
	b.Add(c)
	c.Add(f)
	d.Add(e)
	e.Add(g)

	// Adding d to a should yield:
	//
	// a - b - c - f
	//          `- g
	//
	// with d and e discarded because d have the same Name as b and e has the
	// same Name as c.
	a.Add(d)
	if d.Parent != nil || len(d.Children()) != 0 {
		t.Error("Expected \"d\" to have been discarded.")
	}
	if e.Parent != nil || len(e.Children()) != 0 {
		t.Error("Expected \"e\" to have been discarded.")
	}
	if len(c.Children()) != 2 {
		t.Error("Expected \"c\" to have two Children.")
	}
	if child := c.Child(f.Name); child != f {
		t.Error("Expected \"f\" to be a child of \"c\".")
	}
	if child := c.Child(g.Name); child != g {
		t.Error("Expected \"g\" to be a child of \"c\".")
	}
}

func TestFullPathName(t *testing.T) {
	a := Node{Name: "a"}
	b := Node{Name: "b"}
	c := Node{Name: "c"}
	a.Add(&b)
	b.Add(&c)

	s := c.FullPathName("Q")
	if s != "aQbQc" {
		t.Error("Expected the full path to be aQbQc, but it was", s)
	}

	b.Name = "" // Empty name. Path stops here.
	s = c.FullPathName("Q")
	if s != "c" {
		t.Error("Expected the full path to be c, but it was", s)
	}
}

func TestTraverse(t *testing.T) {
	a := &Node{Name: "A", Kind:KindType}
	b := &Node{Name: "B", Kind:KindNone}
	c := &Node{Name: "C", Kind:KindNone}
	d := &Node{Name: "D", Kind:KindType}
	a.Add(b)
	b.Add(c)
	b.Add(d)

	b.Traverse(func (n *Node) bool {
		return n.Kind == KindType
	}, func (n *Node){
		n.TypeName = "_"
	})
	if a.TypeName != "" {
		t.Error("Did not expect a.TypeName to change")
	}
	if b.TypeName != "" {
		t.Error("Did not expect b.TypeName to change")
	}
	if c.TypeName != "" {
		t.Error("Did not expect c.TypeName to change")
	}
	if d.TypeName != "_" {
		t.Error("Did expect a.TypeName to change to _, but it was", d.TypeName)
	}
}

func TestDumpWithIndent(t *testing.T) {
	a := Node{Name:"A_NAME"}
	b := Node{Name:"B_NAME"}
	a.Add(&b)

	var buf bytes.Buffer
	a.Dump(&buf)
	if buf.Len() == 0 {
		t.Error("Expected the buffer to have some content.")
	}
	if !strings.Contains(buf.String(), "A_NAME") {
		t.Error("Expected the buffer to contain the name of top Node.")
	}
	if !strings.Contains(buf.String(), "B_NAME") {
		t.Error("Expected the buffer to contain the name of second Node.")
	}
}

func TestFind(t *testing.T) {
	a := &Node{Name: "A"}
	b := &Node{Name: "B"}
	c := &Node{Name: "C"}
	d := &Node{Name: "D"}
	a.Add(b)
	b.Add(c, d)

	if n := a.Find("D"); n != d {
		t.Error("Excpected to find d, but got", n)
	}
}

func TestTreeEqual(t *testing.T) {
	a1 := &Node{Name: "A"}
	b1 := &Node{Name: "B"}
	a1.Add(b1)
	a2 := &Node{Name: "A"}
	b2 := &Node{Name: "B"}
	a2.Add(b2)

	if a1.Equal(nil) {
		t.Error("No node can be equal to nil.")
	}
	if !a1.Equal(a2) {
		t.Error("Expected a1 and a2 to be equal.")
	}

	b2.TypeName = "BT"
	if a1.Equal(a2) {
		t.Error("Did not expect a1 and a2 to be equal.")
	}

	b1.TypeName = "BT"
	if !a1.Equal(a2) {
		t.Error("Expected a1 and a2 to be equal.")
	}

	c1 := &Node{Name: "C"}
	b1.Add(c1)
	if a1.Equal(a2) {
		t.Error("Did not expect a1 and a2 to be equal.")
	}
}

func TestGetTop(t *testing.T) {
	a := &Node{Name: "A"}
	b := &Node{Name: "B"}
	c := &Node{Name: "C"}
	d := &Node{Name: "D"}
	a.Add(b)
	b.Add(c, d)

	if top := a.GetTop(); top != a {
		t.Error("Expected a, got", top)
	}
	if top := d.GetTop(); top != a {
		t.Error("Expected a, got", top)
	}
}