package mib

import "fmt"

type DataType struct {
	Name         string
	BaseType     string
	Restrictions string
	CodingType   string
	CodingValue  int
}

func (d DataType) String() string {
	return fmt.Sprintf("DataType: [Name %v, BaseType %v, Restrictions %v, CodingType %v, CodingValue %v.]\n",
		d.Name, d.BaseType, d.Restrictions, d.CodingType, d.CodingValue)
}

type dataTypes []DataType
