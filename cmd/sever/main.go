package main

import (
	"encoding/binary"
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
		//fmt.Println(" ")
		buffer := make([]byte, 1024)
		_, err := c.Read(buffer)
		if err != nil {
			fmt.Println("err", err, "client disconnect")
			fmt.Println("errbuf->", buffer[:50])
			return
		}
		packetId := buffer[0]
		//fmt.Println("GOT", packet.IDToName[packet.ID(packetId)])
		//fmt.Println("\tbuffer->", buffer[:50])
		switch packet.ID(packetId) {
		case packet.KeepAlive:
		case packet.Handshake:
			p := packet.New(
				packet.WithID(packet.Handshake),
				packet.WithString16("-"),
			)

			c.Write(p.Body)

			//fmt.Println("wrote bytes")
		case packet.LoginRequest:
			p := packet.New(
				packet.WithID(packet.LoginRequest),
				packet.WithInt4(1),
				packet.WithString16(""),
				packet.WithLong(22),
				packet.WithByte(0),
			)
			c.Write(p.Body)
		case packet.PlayerPosition:
			//p := packet.New(
			//	packet.WithID(packet.PlayerPosition),
			//	packet.WithDouble(102.809), // X
			//	packet.WithDouble(70.00),   // Y
			//	packet.WithDouble(71.62),   // stance [0.1, 1.65]
			//	packet.WithDouble(68.30),
			//	packet.WithBool(true),
			//)
			//c.Write(p.Body)
		case packet.PlayerLook:
		case packet.PlayerPositionAndLook:
			// 42 b total [OK, get bad packet id 67]
			p := packet.New(
				packet.WithID(packet.PlayerPositionAndLook),
				packet.WithDouble(6.5),   // X
				packet.WithDouble(67.42), // stance
				packet.WithDouble(65.62), // y
				packet.WithDouble(7.5),   // z
				packet.WithDouble(0),     // yaw
				packet.WithDouble(0),     // pitch
				packet.WithBool(false),   // on ground (had to be false to not insta connect)
			)
			//fmt.Println("write->", p.Body)
			c.Write(p.Body)
		case packet.ChatMessage:
			// I probably need to track clients connected here to know the names
			fmt.Println("chat:", string(buffer)) // TODO: parse properly
			fmt.Println(buffer[:50])
			size := binary.BigEndian.Uint16(buffer[1:3])
			message := buffer[3 : size*2+3]
			fmt.Println("size:", size)
			fmt.Println("message:", message)

			p := packet.New(
				packet.WithID(packet.ChatMessage),
				packet.WithString16("<_4dam> "+string(message)),
			)
			c.Write(p.Body)

		default:
			//fmt.Println("unhandled packet, id:", packetId, packet.IDToName[packet.ID(packetId)])
			//fmt.Printf("hex:%x\n", packetId)
			//panic("ooopsie")
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
