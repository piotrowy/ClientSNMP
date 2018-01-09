package mib

type Object interface {
	name() string
	class() string
	number() int
	repr() string
	String() string
}