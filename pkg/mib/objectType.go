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

type objectTypes []ObjectType
