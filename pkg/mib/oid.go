package mib

import (
	"fmt"
)

type Oid struct {
	Value  string
	Name   string
	Class  string
	Number int
}

func (o Oid) String() string {
	return fmt.Sprintf("Oid: [Name %v, Class %v, Number %v.]\n", o.Name, o.Class, o.Number)
}

type oids []Oid

func root(o oids) (oids, Oid, error) {
	for i, v := range o {
		if o.all(func(id Oid) bool {
			return id.Class != v.Name
		}) {
			return append(o[:i], o[i+1:]...), v, nil
		}
	}
	return oids{}, Oid{}, fmt.Errorf("root oid does not exist")
}

func next(o oids) (oids, Oid, error) {
	for i, v := range o {
		if v.Class == v.Name {
			return append(o[:i], o[i+1:]...), v, nil
		}
	}
	return oids{}, Oid{}, fmt.Errorf("oid does not exist")
}

func (o oids) all(fn func(id Oid) bool) bool {
	for _, v := range o {
		if !fn(v) {
			return false
		}
	}
	return true
}
