// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Client Package provides a TCP client for communicating with the dataserver.
package client

import (
	"net"

	"github.com/teonet-go/dataserver"
)

const (
	// startPacketLength is the length in bytes of the start packet
	// sent by the client to initialize the connection.
	// reserved for the future constant:
	// startPacketLength = dataserver.StartPacketLength

	// ChankPacketLength is the max length in bytes of each data
	// packet sent over the connection.
	ChankPacketLength = 21

	// timeout is the read/write timeout for network operations.
	// reserved for the future constant:
	// timeout = 5 * time.Second
)

// DataClient is the client data structure containing the remote address,
// start packet, and network connection.
type DataClient struct {
	remoteAddr  string
	startPacket *dataserver.StartPacket
	net.Conn
}

// NewDataClient creates a new DataClient instance, connects to the DataServer
// specified by remoteAddr, and sends the provided startPacket to initialize
// the connection. It returns the new DataClient and any error.
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

// Read reads data from the connection into the provided byte slice b.
// It returns the number of bytes read and any error encountered.
// Read can be made to time out and return an error after a fixed
// time limit by setting deadlines on the connection.
func (dc *DataClient) Read(b []byte) (n int, err error) {
	return dc.Conn.Read(b)
}
