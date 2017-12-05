package mib

import "testing"

const (
	VALUE  = "internet.mib"
	NAME   = "mib"
	PARENT = "internet"
)

var (
	oid    = NewOid(VALUE, NAME, "", 0)
	oidParent = NewOid(PARENT, PARENT,"", 0)
)

func TestParentOid(t *testing.T) {
	parent := oid.Parent()
	if parent.Name != PARENT && parent.Value == PARENT {
		t.Errorf("Parent oid is different than expected")
	}
}

func TestMatchOid(t *testing.T) {
	if _, ok := oid.Match(oidParent); ok != nil {
		t.Errorf("Matching wet wrong.")
	}
}

func TestId(t *testing.T) {
	if oid.Id() != NAME {
		t.Errorf("Id is different than expected")
	}
}
