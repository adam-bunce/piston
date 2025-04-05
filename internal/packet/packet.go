package packet

type Packet interface {
	Type()
	Size() uint64
	Serialize() []byte
}

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

//  i need a way to handle the creation of packets easily
// so might as well make builder to abstract away difficulty of
// doing little endian utf16 conversion
// with string
// with long
// with int
// with byte
// with bytes([]byte) that just gets dumped
// with normalBytes (we do conversion to utf16 or whateva? )

// GOAL 1: two users connect, can send chat messages (no blocks or anything else)
