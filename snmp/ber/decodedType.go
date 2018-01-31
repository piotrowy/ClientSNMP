package ber

const (
	UNIVERSAL = iota
	APPLICATION
	CONTEXT_SPECIFIC
	PRIVATE
)

type decodedType struct {
	typeTagId int
	visibilityClass int
	data []byte
	length uint32
	isConstructed bool
}
