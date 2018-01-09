package mib

import (
	"fmt"
	"sort"
	"bytes"

	"strings"
)

type (
	Tree struct {
		root *node
		size int
	}

	WalkFn func(n *node) bool
)

func New(oid Oid, val ObjectType) *Tree {
	return &Tree{
		root: &node{
			parent:   nil,
			children: []*node{},
			id:       oid,
			val:      val,
			height:   0,
		},
		size: 1,
	}
}

func (t *Tree) Len() int {
	return t.size
}

func (t *Tree) Walk(fn WalkFn) {
	recursiveDfsWalk(t.root, fn)
}

func (t *Tree) InsertOid(id Oid) {
	t.Insert(id, ObjectType{})
}

func (t *Tree) Insert(id Oid, val ObjectType) {
	var (
		n   *node
		err error
	)
	n, err = t.root.findByName(id.Class)
	if err == nil {
		n.insert(id, val)
		t.size += 1
	}
}

func (t *Tree) Delete(id Oid) {
	n, err := t.root.findByName(id.Class)
	if err != nil {
		n.delete(id)
		t.size -= 1
	}
}

func (t *Tree) FindByOid(id Oid) (Oid, ObjectType, error) {
	n, err := t.root.findByOid(id)
	return n.id, n.val, err
}

func (t *Tree) FindByName(name string) (Oid, ObjectType, error) {
	n, err := t.root.findByName(name)
	return n.id, n.val, err
}

func (t *Tree) String() string {
	return t.root.string(Oid{})
}

func (t *Tree) SubtreeString(id Oid) string {
	if _, _, err := t.FindByOid(id); err != nil {
		panic(err)
	}
	return t.root.string(id)
}

func (t *Tree) toMap() map[string]ObjectType {
	out := make(map[string]ObjectType)
	t.Walk(func(n *node) bool {
		out[n.id.Value] = n.val
		return false
	})
	return out
}

type node struct {
	parent   *node
	children children
	id       Oid
	val      ObjectType
	height   int
}

func (n *node) insert(id Oid, val ObjectType) {
	if n.indexOf(id) == -1 {
		n.children = append(n.children, &node{
			parent:   n,
			children: []*node{},
			id:       id,
			val:      val,
			height:   n.height + 1,
		})
		n.children.Sort()
	}
}

func (n *node) delete(id Oid) {
	i := n.indexOf(id)
	if i >= 0 {
		n.children = append(n.children[:i], n.children[i+1:]...)
	}
}

func (n *node) findByName(name string) (*node, error) {
	return n.findBy(func(id Oid) bool {
		return id.Name == name
	})
}

func (n *node) findByOid(id Oid) (*node, error) {
	return n.findBy(func(id2 Oid) bool {
		return id2 == id
	})
}

func (n *node) findBy(fn func(id Oid) bool) (*node, error) {
	var res *node
	ok := recursiveDfsWalk(n, func(node *node) bool {
		if fn(node.id) {
			res = node
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

func (n *node) string(id Oid) string {
	var buff bytes.Buffer
	bfsWalk(n, func(node *node) bool {
		str := fmt.Sprintf("%v=> %v\n", strings.Repeat(space, node.height*4), node.id.Name)
		buff.WriteString(str)
		if node.id.Name == id.Name {
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
		b.WriteString(fmt.Sprintf("=> %v", c.id.Name))
	}
}

func (n *node) indexOf(id Oid) int {
	for i, v := range n.children {
		if v.id == id {
			return i
		}
	}
	return -1
}

func (n node) name() (r string) {
	if n.id != (Oid{}) {
		r = fmt.Sprintf("{%v}: {%v}", n.id.Number, n.id.Name)
	} else {
		r = fmt.Sprintf("{%v}: {%v}", n.val.Number, n.val.Name)
	}
	return
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
	return c[i].id.Number < c[j].id.Number
}

func (c children) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c children) Sort() {
	sort.Sort(c)
}
