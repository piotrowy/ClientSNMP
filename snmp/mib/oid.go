package mib

import (
	"fmt"
)

type Oid struct {
	Name   string
	Class  string
	Number int
}

func (o Oid) GetName() string {
	return o.Name
}

func (o Oid) GetClass() string {
	return o.Class
}

func (o Oid) GetNumber() int {
	return o.Number
}

func (o Oid) Representation() string {
	return fmt.Sprintf("{%v}: {%v}", o.Number, o.Name)
}

func (o Oid) String() string {
	return fmt.Sprintf("Oid: [Name %v, Class %v, Number %v]\n", o.Name, o.Class, o.Number)
}

type oids []Oid

func (o oids) next() (oids, Oid, error) {
	for i, v := range o {
		if o.all(func(id Oid) bool {
			return id.Name != v.Class
		}) {
			return append(o[:i], o[i+1:]...), v, nil
		}
	}
	return oids{}, Oid{}, fmt.Errorf("no next oid")
}

func (o oids) all(fn func(id Oid) bool) bool {
	for _, v := range o {
		if !fn(v) {
			return false
		}
	}
	return true
}
