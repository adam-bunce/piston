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
			c.Write([]byte{2, 0, 1, 0, '-'})
			fmt.Println("wrote bytes")
		case packet.LoginRequest:
			// player id 4 bytes
			// Unknown string 16 (wtf does this mean ? always empty?)
			// map seed (not used by client, shouldnt send if i gaf about sec)
			// dimension // (
			// 4 + 2 + 8 + 1
			//                                 |   player id(4)    |  unk(2 + strlen) | map seed (8)             | dim
			c.Write([]byte{packet.LoginRequest, 0, 0, 0, 1, 0, 0, 1, 2, 3, 5, 5, 6, 7, 8, 0})
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
