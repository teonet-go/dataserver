package client

import (
	"net"
	"time"

	"github.com/teonet-go/dataserver"
)

const (
	startPacketLength = dataserver.StartPacketLength
	ChankPacketLength = 21
	timeout           = 5 * time.Second
)

// DataClient is TCP Data Client data structure and methods receiver.
type DataClient struct {
	remoteAddr  string
	startPacket *dataserver.StartPacket
	net.Conn
}

// NewDataClient creates new DataClient, connects to DataServer and sends
// read or write request depending on the type of start packet.
func NewDataClient(remoteAddr string, startPacket *dataserver.StartPacket) (
	dc *DataClient, err error) {

	// Create new DataClient object
	dc = &DataClient{remoteAddr: remoteAddr, startPacket: startPacket}

	// Connect to the server
	dc.Conn, err = net.Dial("tcp", remoteAddr)
	if err != nil {
		return
	}

	// Send start packet to the server
	_, err = dc.Write(dc.startPacket.Bytes())
	if err != nil {
		dc.Close()
		return
	}

	return
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (dc *DataClient) Write(b []byte) (n int, err error) {
	return dc.Conn.Write(b)
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (dc *DataClient) Read(b []byte) (n int, err error) {
	return dc.Conn.Read(b)
}
