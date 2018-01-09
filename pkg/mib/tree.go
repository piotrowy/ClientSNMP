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

func New(oid Oid, val ObjectType) *Tree {
	return &Tree{
		root: &node{
			parent:   nil,
			children: []*node{},
			id:       oid,
			val:      val,
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

func (t *Tree) InsertOid(id Oid) {
	t.insert(id, ObjectType{})
}

func (t *Tree) InsertObjectType(ot ObjectType) {
	t.insert(Oid{}, ot)
}

func (t *Tree) insert(id Oid, val ObjectType) {
	n, err := t.root.findByName(node{
		id:  id,
		val: val,
	}.class())
	if err == nil {
		n.insert(id, val)
		t.size += 1
	}
}

func (t *Tree) Delete(id Oid, val ObjectType) {
	n2 := node{
		id:  id,
		val: val,
	}
	n, err := t.root.findByName(n2.class())
	if err != nil {
		n.delete(n2)
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
	return t.root.string(Oid{}, ObjectType{})
}

func (t *Tree) SubtreeString(id Oid, ot ObjectType) string {
	if _, _, err := t.FindByName(node{
		id:  id,
		val: ot,
	}.name()); err != nil {
		panic(err)
	}
	return t.root.string(id, ot)
}

type node struct {
	parent   *node
	children children
	id       Oid
	val      ObjectType
	height   int
}

func (n *node) insert(id Oid, val ObjectType) {
	if n.indexOf(node{
		id:  id,
		val: val,
	}) == -1 {
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

func (n *node) delete(node node) {
	i := n.indexOf(node)
	if i >= 0 {
		n.children = append(n.children[:i], n.children[i+1:]...)
	}
}

func (n *node) findByName(name string) (*node, error) {
	return n.findBy(func(node node) bool {
		return node.name() == name
	})
}

func (n *node) findByOid(id Oid) (*node, error) {
	return n.findBy(func(node node) bool {
		return node.id == id
	})
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

func (n *node) string(id Oid, ot ObjectType) string {
	var buff bytes.Buffer
	bfsWalk(n, func(n2 *node) bool {
		str := fmt.Sprintf("%v=> %v\n", strings.Repeat(space, n2.height*4), n2.repr())
		buff.WriteString(str)
		if n2.name() == (node{
			id:  id,
			val: ot,
		}.name()) {
			n2.writeChildren(buff)
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
		b.WriteString(fmt.Sprintf("=> %v", c.name()))
	}
}

func (n *node) indexOf(node node) int {
	for i, v := range n.children {
		if v.name() == node.name() {
			return i
		}
	}
	return -1
}

func (n node) name() (r string) {
	if n.id != (Oid{}) {
		r = n.id.Name
	} else {
		r = n.val.Name
	}
	return
}

func (n node) class() (r string) {
	if n.id != (Oid{}) {
		r = n.id.Class
	} else {
		r = n.val.Class
	}
	return
}

func (n node) repr() (r string) {
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
