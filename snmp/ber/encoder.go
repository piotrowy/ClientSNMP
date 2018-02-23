package ber

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"strings"
)

const (
	sequence    = "SEQUENCE"
	integer     = "INTEGER"
	octetString = "STRING"
)

var (
	TagNumbers = map[string]int{
		"INTEGER":           2,
		"BIT STRING":        3,
		"OCTET STRING":      4,
		"NULL":              5,
		"OBJECT IDENTIFIER": 6,
		"SEQUENCE":          16,
	}
)

func Encode(syntax, value string) (res []byte) {
	syntax = clearSyntax(syntax)
	res = append(res, encodeTag(syntax)...)
	switch syntax {
	case integer:
		if i, err := strconv.ParseInt(value, 10, 16); err == nil {
			res = append(res, encodeInt(int(i))...)
		}
	case octetString:
	case sequence:
	}

	return
}

func encodeTag(syntax string) []byte {
	tag := 0
	if isConstructed(syntax) {
		tag |= 1 << 5
	}
	if v, ok := TagNumbers[syntax]; ok {
		tag |= v
	}
	return []byte{byte(uint16(tag))}
}

func encodeInt(i int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, clearInt(i))
	if err != nil {
		panic(err)
	}
	b := buf.Bytes()
	return append(encodeLength(len(b)), b...)
}

func encodeBool(b bool) []byte {
	return []byte{}
}

func encodeNull() []byte {
	return []byte{}
}

func encodeLength(len int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, clearInt(len))
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func reverse(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

func encodeString() []byte {
	return []byte{}
}

func isConstructed(syntax string) bool {
	return syntax == sequence
}

func clearSyntax(syntax string) (res string) {
	if strings.Contains(syntax, integer) {
		res = integer
	} else if strings.Contains(syntax, octetString) {
		res = octetString
	} else if strings.Contains(syntax, sequence) {
		res = sequence
	} else {
		res = syntax
	}
	return
}

func clearInt(val int) (res interface{}) {
	if math.MinInt8 < val && val < math.MaxInt8 {
		res = int8(val)
	} else if math.MinInt16 < val && val < math.MaxInt16 {
		res = int16(val)
	} else if math.MinInt32 < val && val < math.MaxInt32 {
		res = int32(val)
	} else {
		res = int64(val)
	}
	return
}
