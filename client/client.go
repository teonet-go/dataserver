package client

import (
	"net"
	"time"

	"github.com/teonet-go/teoss/dataserver"
)

const (
	startPacketLength = 6
	ChankPacketLength = 21
	timeout           = 5 * time.Second
)

// DataClient is TCP Data Client data structure and methods receiver.
type DataClient struct {
	remoteAddr string
	request    dataserver.StartPacket
	net.Conn
}

// NewDataClientWriter creates new DataClient, connects to DataServer and sends
// request with start packets data.
func NewDataClientWriter(remoteAddr string, request dataserver.StartPacket) (
	dc *DataClient, err error) {

	// Create new DataClient object
	dc = &DataClient{remoteAddr: remoteAddr, request: request}

	// Connect to the server
	dc.Conn, err = net.Dial("tcp", remoteAddr)
	if err != nil {
		return
	}

	// Send start packet to the server
	_, err = dc.Write(dc.request.Bytes())
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
