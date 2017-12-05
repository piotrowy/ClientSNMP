package mib

import (
	"strings"
	"regexp"
)

const separator = "."

type Oid struct {
	Value  string
	Name   string
	Class  string
	Number int
}

func (oid *Oid) Parent() (o Oid) {
	arr := strings.Split(oid.Value, separator)
	o.Name = arr[len(arr)-2]
	for _, v := range arr[:len(arr)-1] {
		o.Value += v
	}
	return
}

func (oid *Oid) Match(o Oid) (bool, error) {
	return regexp.MatchString(oid.Value, o.Value)
}

func (oid *Oid) String() string {
	return oid.Value
}

func NewOid(oid, name, class string, number int) (o Oid) {
	o = Oid{
		Value:  oid,
		Name:   name,
		Class:  class,
		Number: number,
	}
	return
}

func ShortOid(oid string) (o Oid) {
	slices := strings.Split(oid, separator)
	o = Oid{
		Value: oid,
		Name: slices[len(slices)-1],
	}
	return
}
