package oid

import "testing"

const (
	VALUE  = "internet.mib"
	NAME   = "mib"
	PARENT = "internet"
)

var (
	oid    = New(VALUE, NAME)
	oidParent = New(PARENT, PARENT)
)

func TestParentOid(t *testing.T) {
	parent := oid.ParentOid()
	if parent.Name != PARENT && parent.Value == PARENT {
		t.Errorf("Parent oid is different than expected")
	}
}

func TestMatchOid(t *testing.T) {
	if _, ok := oid.MatchOid(oidParent); ok != nil {
		t.Errorf("Matching wet wrong.")
	}
}

func TestId(t *testing.T) {
	if oid.Id() != NAME {
		t.Errorf("Id is different than expected")
	}
}
