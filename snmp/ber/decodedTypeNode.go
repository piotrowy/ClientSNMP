package ber

type decodedTypeNode struct {
	parent   *decodedTypeNode
	children []*decodedTypeNode
	value    decodedType
}
