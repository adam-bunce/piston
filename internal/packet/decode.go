package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
)

//type ServerBound struct {
//	ID
//	Values map[string]interface{}
//}

type Field int

const (
	unknownF Field = iota
	doubleF
	floatF

	intF
	longF

	boolF
	string16F
	byteF
)

// interface here maybe

type Structure struct {
	Name string
	Field
	Optional bool
}

var sizes = map[Field]int{
	doubleF: 8,
	floatF:  4,

	longF: 8,
	intF:  4,

	string16F: 2, // min 2

	boolF: 1,
	byteF: 1,
}

var serverBoundPackets = map[ID][]Structure{
	KeepAlive: {},
	Handshake: {
		{"username", string16F, false},
	},
	LoginRequest: {
		{"protocolVersion", intF, false},
		{"username", string16F, false},
		//{"mapSeed", longF, true},
		//{"dimension", byteF, true},
	},
	PlayerPositionAndLook: {
		{"x", doubleF, false},
		{"y", doubleF, false},
		{"stance", doubleF, false},
		{"z", doubleF, false},

		{"yaw", floatF, false},
		{"pitch", floatF, false},
		{"onGround", boolF, false}, // derived from 'flying' packet not read so maybe don't include here or make implicit
	},
	ChatMessage: {
		{"message", string16F, false},
	},
}

func GetPacketType(c net.Conn) ID {
	buffer := make([]byte, 1)
	// assert we read 1 by and that byte exists in packet id's
	read, err := c.Read(buffer)
	if err != nil {
		fmt.Println("FAILED TO READ FROM BUFFER DURING PACKET GET")
		return 0
	}
	fmt.Println("read", read, "bytes", IDToName[ID(buffer[0])])
	return ID(buffer[0])
}

// ParsePacket feel like implementation mwill matter dependino how how i used the data (map is okayish ig)
func ParsePacket(c net.Conn) (map[string]any, error) {
	pt := GetPacketType(c)

	parts := serverBoundPackets[pt]
	packetParts := make(map[string]any)
	packetParts["id"] = pt

	for i, part := range parts {
		size := sizes[part.Field]
		buffer := make([]byte, size)

		_, err := io.ReadFull(c, buffer)
		if err != nil {
			return nil, err
		}

		// type cast map Field->func?
		switch part.Field {
		case intF:
			packetParts[part.Name] = int(binary.BigEndian.Uint32(buffer))
		case string16F:
			stringSize := binary.BigEndian.Uint16(buffer[0:2])
			packetString := make([]byte, stringSize*2)
			_, err = io.ReadFull(c, packetString)
			if err != nil {
				return nil, err
			}
			packetParts[part.Name] = UTF16ToRunes(packetString)
		case longF:
			packetParts[part.Name] = LongToInt(buffer)
		case byteF:
			packetParts[part.Name] = buffer
		case doubleF:
			packetParts[part.Name] = math.Float64frombits(binary.BigEndian.Uint64(buffer))
		case floatF:
			//packetParts[part.Name] = buffer
			packetParts[part.Name] = math.Float32frombits(binary.BigEndian.Uint32(buffer))
		case boolF:
			packetParts[part.Name] = buffer[0] == 1
		default:
			panic("unhandled packet part case")
		}

		fmt.Printf("%d %q %v\n", i, part.Name, packetParts[part.Name])
	}

	return packetParts, nil

}
