package mib

import (
	"fmt"
)

type ObjectType struct {
	Name        string
	Syntax      string
	Access      string
	Status      string
	Description string
	Class       string
	Number      int
	Min         uint64
	Max         uint64
}

func (o ObjectType) String() string {
	return fmt.Sprintf("ObjectType: [Name %v, Syntax %v, Access %v, "+
		"Status %v, Description %v, Class %v, Number %v, Min %v, Max %v.]\n",
		o.Name, o.Syntax, o.Access, o.Status, o.Description, o.Class, o.Number,
		o.Min, o.Max)
}

func (o ObjectType) name() string {
	return o.Name
}

func (o ObjectType) class() string {
	return o.Class
}

func (o ObjectType) number() int {
	return o.Number
}

func (o ObjectType) repr() string {
	return fmt.Sprintf("{%v}: {%v}", o.Number, o.Name)
}

type objectTypes []ObjectType

func (ots objectTypes) next() (objectTypes, ObjectType, error) {
	if len(ots) > 0 {
		return ots[1:], ots[0], nil
	}
	return objectTypes{}, ObjectType{}, fmt.Errorf("no next type")
}