package packet

import "encoding/binary"

//type Packet interface {
//	Type()
//	Size() uint64
//	Serialize() []byte
//}

/*
for server bound packet parsing I need to read
the first byte (packet_id), then hand the byte stream off to that particular parser
*/

type ID byte

const (
	KeepAlive             ID = 0x0
	LoginRequest             = 0x01
	Handshake                = 0x02
	PlayerPosition           = 0x0B
	PlayerLook               = 0x0C
	PlayerPositionAndLook    = 0x0D
)

var IDToName = map[ID]string{
	KeepAlive:             "Keep Alive",
	LoginRequest:          "Login Request",
	Handshake:             "Handshake",
	PlayerPosition:        "Player Position",
	PlayerLook:            "Player Look",
	PlayerPositionAndLook: "Player Position and Look",
}

// build/options pattern for creating packets
/*
p := packet.newPacket(
	packet.withId(packet.Handshake)
	packet.withString("_4dam") -> adds uint16, little endian to end of packet position
)

*/

// pretty print packet to do decode and encode tests

// something similar to 8086_sim to convert bytes to usable objects would be good
// except packet id give it away so easier
// describe packet shape, get out packet struct (or map? idk)

// GOAL 1: two users connect, can send chat messages (no blocks or anything else)

type Packet struct {
	Body []byte
}

func New(options ...func(*Packet)) *Packet {
	packet := &Packet{Body: []byte{}}
	for _, function := range options {
		function(packet)
	}

	return packet
}

// WithByte writes the given byte end of the packet
func WithByte(b byte) func(*Packet) {
	return func(p *Packet) {
		p.Body = append(p.Body, b)
	}
}

// WithID writes the given byte to the first position in the packet
func WithID(id byte) func(*Packet) {
	return func(p *Packet) {
		if len(p.Body) == 0 {
			p.Body = make([]byte, 1)
		}
		p.Body[0] = id
	}
}

// WithString16 writes s as 2 + len(s) bytes to the packet
func WithString16(s string) func(*Packet) {
	return func(p *Packet) {
		p.Body = append(p.Body, StringToBytes(s)...)
	}
}

// WithInt4 writes v as 4 bytes to the packet
func WithInt4(v int) func(*Packet) {
	return func(p *Packet) {
		p.Body = binary.BigEndian.AppendUint32(p.Body, uint32(v))
	}
}

// WithLong writes v as 8 bytes to the packet
func WithLong(v int) func(*Packet) {
	return func(p *Packet) {
		//p.Body = append(p.Body, 0, 0, 0, 0)
		p.Body = binary.BigEndian.AppendUint64(p.Body, uint64(v))
	}
}
