package main

import (
	"fmt"
	"net"
	"os"
	"piston/internal/packet"
)

func errorExit(e error, s ...string) {
	if e != nil {
		fmt.Println(s, e)
		os.Exit(1)
	}
}

func handleClient(c net.Conn) {
	fmt.Println("client connected from:", c.RemoteAddr())
	defer c.Close()

	for {
		packetData, _ := packet.ParsePacket(c)
		fmt.Println("packet data", packetData)

		// packet data map[id:2 username:[95 52 100 97 109]]

		switch packetData["id"].(packet.ID) {
		case packet.KeepAlive:
		case packet.Handshake:
			p := packet.New(
				packet.WithID(packet.Handshake),
				packet.WithString16("-"),
			)

			c.Write(p.Body)
		case packet.LoginRequest:
			const (
				EntityId  = 1
				MapSeed   = 12345
				Dimension = 0
			)
			p := packet.New(
				packet.WithID(packet.LoginRequest),
				packet.WithInt4(EntityId),
				packet.WithString16("piston"),
				packet.WithLong(MapSeed),
				packet.WithByte(Dimension),
			)
			c.Write(p.Body)
		case packet.PlayerPosition:
		case packet.PlayerLook:
		case packet.PlayerPositionAndLook:
			// 42 b total [OK, get bad packet id 67]
			p := packet.New(
				packet.WithID(packet.PlayerPositionAndLook),
				packet.WithDouble(6.5),                         // X
				packet.WithDouble(67.42),                       // stance
				packet.WithDouble(65.62),                       // y
				packet.WithDouble(7.5),                         // z
				packet.WithDouble(0),                           // yaw
				packet.WithDouble(0),                           // pitch
				packet.WithBool(packetData["onGround"].(bool)), // on ground (had to be false to not insta connect)
			)
			c.Write(p.Body)
		case packet.ChatMessage:
			// I probably need to track clients connected here to know the names
			p := packet.New(
				packet.WithID(packet.ChatMessage),
				packet.WithString16("<_4dam> "+string(packetData["message"].([]rune))),
			)
			c.Write(p.Body)
		default:
			fmt.Println("unknown packet")
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", "localhost:25565")
	errorExit(err, "net.Listen")

	// listen for client connection attempts
	for {
		conn, err := ln.Accept()
		errorExit(err)
		fmt.Println("made conn")

		go handleClient(conn)
	}
}
