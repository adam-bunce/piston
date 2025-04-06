package packet

import (
	"encoding/binary"
)

// StringToBytes converts the given string into UTF-16 bytes with a 2 byte length prefix
func StringToBytes(s string) []byte {
	var buffer []byte

	buffer = binary.BigEndian.AppendUint16(buffer, uint16(len(s)))

	for _, char := range s {
		buffer = binary.BigEndian.AppendUint16(buffer, uint16(char))
	}

	return buffer
}
