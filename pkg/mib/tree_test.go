package mib

import (
	"testing"
)

var tr = New(Oid{
	Name:   "internet",
	Class:  "iso",
	Number: 1,
}, ObjectType{})

func TestInsert(t *testing.T) {
	//when
	tr.InsertOid(Oid{
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
