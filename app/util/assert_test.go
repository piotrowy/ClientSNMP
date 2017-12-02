package util

import "testing"

func TestAssertPanic(t *testing.T) {
	defer errorWhenNoPanic(t)
	Assert(false)
}

func TestAssertPPanic(t *testing.T) {
	defer errorWhenNoPanic(t)
	AssertP(func() bool { return false })
}

func TestAssert(t *testing.T) {
	defer errorWhenPanic(t)
	Assert(true)
}

func TestAssertP(t *testing.T) {
	defer errorWhenPanic(t)
	AssertP(func() bool { return true })
}

func errorWhenNoPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}

func errorWhenPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Errorf("The code should not panic")
	}
}
