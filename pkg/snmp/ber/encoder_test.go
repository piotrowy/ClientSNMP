package ber

import (
	"testing"
	"bytes"
)

func TestEncodeInteger(t *testing.T) {
	encoded := Encode(integer, "-129")
	expected := []byte{0x02, 0x02, 0xFF, 0x7F}
	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encoding integer failed. Expected <%v>\nGot <%v>", expected, encoded)
	}
}

func TestEncodeNull(t *testing.T) {
	encoded := Encode(integer, "NULL")
	expected := []byte{0x05, 0x00}
	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encoding null failed. Expected <%v>\nGot <%v>", expected, encoded)
	}
}