package mib

import (
	"testing"
)

func given() *Tree {
	return New(Oid{
		Name:   "internet",
		Class:  "iso",
		Number: 1,
	})
}

func TestInsertOid(t *testing.T) {
	//when
	tr := given()
	tr.Insert(Oid{
		Name:   "directory",
		Class:  "internet",
		Number: 1,
	})

	//then
	if len(tr.root.children) <= 0 {
		t.Errorf("Inserting failed.")
	}

	//after
	tr.root.children = []*node{}
}
