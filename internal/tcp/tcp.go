// Package tcp handles client-server connections and the transport of data packets
package tcp

const defaultPort = 25565

type Connection interface {
	Send([]byte)
	// Receive using context here might be good, if someone's packet
	// cant be handled in n seconds then give up and cancel all child processes
	Receive() []byte
}

type ClientServerConnection struct {
	clientIp   string
	serverIp   string
	clientPort int
	serverPort int
}
