package mib

import (
	"testing"
)

var tr = New(Oid{
	Value:  "",
	Name:   "internet",
	Class:  "iso",
	Number: 1,
}, ObjectType{})

func TestInsert(t *testing.T) {
	tr.InsertOid(Oid{
		Value:  "",
		Name:   "directory",
		Class:  "internet",
		Number: 1,
	})
	if len(tr.root.children) <= 0 {
		t.Errorf("Inserting failed.")
	}
}
