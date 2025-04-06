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
		fmt.Println(" ")
		buffer := make([]byte, 1024)
		_, err := c.Read(buffer)
		if err != nil {
			fmt.Println("err", err, "client disconnect")
			fmt.Println("errbuf->", buffer[:50])
			return
		}
		packetId := buffer[0]
		fmt.Println("GOT", packet.IDToName[packet.ID(packetId)])
		fmt.Println("\tbuffer->", buffer[:50])
		switch packet.ID(packetId) {
		case packet.KeepAlive:
		case packet.Handshake:
			p := packet.New(
				packet.WithID(packet.Handshake),
				packet.WithString16("-"),
			)

			c.Write(p.Body)

			fmt.Println("wrote bytes")
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
		case packet.PlayerLook:
		case packet.PlayerPositionAndLook:
		default:
			fmt.Println("unhandled packet, id:", packetId, packet.IDToName[packet.ID(packetId)])
			panic("ooopsie")
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
