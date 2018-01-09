package mib

import (
	"fmt"
	"sort"
	"bytes"

	"strings"
)

type (
	Tree struct {
		root          *node
		size          int
		distinctTypes map[string]interface{}
	}

	WalkFn func(n *node) bool
)

func New(obj Object) *Tree {
	return &Tree{
		root: &node{
			parent:   nil,
			children: []*node{},
			val: obj,
			height:   0,
		},
		size:          1,
		distinctTypes: map[string]interface{}{},
	}
}

func (t *Tree) Len() int {
	return t.size
}

func (t *Tree) Walk(fn WalkFn) {
	recursiveDfsWalk(t.root, fn)
}

func (t *Tree) Insert(o Object) {
	n, err := t.root.findByName(o.class())
	if err == nil {
		n.insert(o)
		t.size += 1
	}
}

func (t *Tree) Delete(o Object) {
	n, err := t.root.findByName(o.class())
	if err != nil {
		n.delete(o)
		t.size -= 1
	}
}

func (t *Tree) FindByOid(id Oid) (Object, error) {
	n, err := t.root.findByObject(id)
	return n.val, err
}

func (t *Tree) FindByName(name string) (Object, error) {
	n, err := t.root.findByName(name)
	return n.val, err
}

func (t *Tree) FindByValue(v string) (Object, error) {
	n, err := t.root.findByValue(v)
	return n.val, err
}

func (t *Tree) String() string {
	return t.root.string(Oid{})
}

func (t *Tree) SubtreeString(o Object) string {
	if _, err := t.FindByName(o.name()); err != nil {
		panic(err)
	}
	return t.root.string(o)
}

type node struct {
	parent   *node
	children children
	val      Object
	height   int
}

func (n *node) insert(val Object) {
	if n.indexOf(val) == -1 {
		n.children = append(n.children, &node{
			parent:   n,
			children: []*node{},
			val:      val,
			height:   n.height + 1,
		})
		n.children.Sort()
	}
}

func (n *node) delete(o Object) {
	i := n.indexOf(o)
	if i >= 0 {
		n.children = append(n.children[:i], n.children[i+1:]...)
	}
}

func (n *node) findByName(name string) (*node, error) {
	return n.findBy(func(node node) bool {
		return node.val.name() == name
	})
}

func (n *node) findByObject(o Object) (*node, error) {
	return n.findBy(func(node node) bool {
		return node.val == o
	})
}

func (n *node) findByValue(v string) (*node, error) {
	return &node{}, nil
}

func (n *node) findBy(fn func(n1 node) bool) (*node, error) {
	var res *node
	ok := recursiveDfsWalk(n, func(n2 *node) bool {
		if fn(*n2) {
			res = n2
			return true
		}
		return false
	})

	if ok {
		return res, nil
	} else {
		return &node{}, fmt.Errorf("cannot find node")
	}
}

func (n *node) string(o Object) string {
	var buff bytes.Buffer
	recursiveDfsWalk(n, func(node *node) bool {
		str := fmt.Sprintf("%v=> %v\n", strings.Repeat(space, node.height*4), node.val.repr())
		buff.WriteString(str)
		if node.val.name() == o.name() {
			node.writeChildren(buff)
			return true
		}
		return false
	})
	return buff.String()
}

func (n *node) writeChildren(b bytes.Buffer) {
	for _, c := range n.children {
		b.WriteString("\n")
		for i := 0; i < c.height; i++ {
			b.WriteString("    ")
		}
		b.WriteString(fmt.Sprintf("=> %v", c.val.name()))
	}
}

func (n *node) indexOf(o Object) int {
	for i, v := range n.children {
		if v.val.name() == o.name() {
			return i
		}
	}
	return -1
}

//recursiveWalk returns true if it should be aborted
func recursiveDfsWalk(n *node, fn WalkFn) bool {
	if fn(n) {
		return true
	}
	for _, n := range n.children {
		if recursiveDfsWalk(n, fn) {
			return true
		}
	}
	return false
}

//bfsWalk returns true if it should be aborted
func bfsWalk(n *node, fn WalkFn) bool {
	nodes := []*node{n}
	for len(nodes) > 0 {
		node := nodes[0]
		nodes = append(nodes[1:], node.children...)
		if fn(node) {
			return true
		}
	}
	return false
}

type children []*node

func (c children) Len() int {
	return len(c)
}

func (c children) Less(i, j int) bool {
	return c[i].val.number() < c[j].val.number()
}

func (c children) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c children) Sort() {
	sort.Sort(c)
}
