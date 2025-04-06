package packet

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

// Encode section

// StringToBytes converts the given string into UTF-16 bytes with a 2 byte length prefix
func StringToBytes(s string) []byte {
	// TODO: assert len >= 2 (always have 2 bytes of length info)
	var buffer []byte

	buffer = binary.BigEndian.AppendUint16(buffer, uint16(len(s)))

	for _, char := range s {
		buffer = binary.BigEndian.AppendUint16(buffer, uint16(char))
	}

	return buffer
}

// Decode section

func UTF16ToRunes(b []byte) []rune {
	// TODO: assert %2 is ok
	var asUint16 []uint16
	for i := 0; i < len(b); i += 2 {
		u := binary.BigEndian.Uint16(b[i : i+2])
		asUint16 = append(asUint16, u)
	}

	return utf16.Decode(asUint16)
}

func LongToInt(b []byte) int {
	// NOTE: not sure if this works
	val := int(binary.BigEndian.Uint64(b))
	fmt.Println(val)
	return val
}
