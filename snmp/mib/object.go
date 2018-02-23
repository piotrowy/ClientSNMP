package mib

type ObjectIdentifier interface {
	GetName() string
	GetClass() string
	GetNumber() int
	Representation() string
	String() string
}