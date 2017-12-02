package oid

import (
	"strings"
	"regexp"
)

type Oid struct {
	Value  string
	Name   string
	Class  string
	Number int
	sep    string
}

func New(oid, name string) (o Oid) {
	o = Oid{
		Value: oid,
		Name:  name,
		sep:   ".",
	}
	return
}

func (oid Oid) ParentOid() (o Oid) {
	arr := strings.Split(oid.Value, oid.sep)
	o.Name = arr[len(arr)-2]
	for _, v := range arr[:len(arr)-1] {
		o.Value += v
	}
	return
}

func (oid Oid) Id() string {
	slices := strings.Split(oid.Value, oid.sep)
	return slices[len(slices)-1]
}

func (oid Oid) MatchOid(o Oid) (bool, error) {
	return regexp.MatchString(oid.Value, o.Value)
}

func (oid Oid) String() string {
	return oid.Value
}
