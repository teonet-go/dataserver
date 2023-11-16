package dataserver

import (
	"crypto/md5"
	"fmt"
)

const (
	READ  = iota // Receive data on data server request.
	WRITE        // Send data from data server request.
)

const StartPacketLength = 17

// StartPacket contains start packet data and type, and methods to murshal &
// unmarshal it.
type StartPacket struct {
	// Type of start packet:
	//  0 - READ (server will receive data);
	//  1 - WRITE (server will send data).
	Type byte

	// Start packet data
	Data []byte
}

// MakeStartPacket creates new start packet from Type and Name
func MakeStartPacket(t byte, name string) (startPacket *StartPacket) {

	startPacket = &StartPacket{Type: t}

	h := md5.New()
	h.Write([]byte(name))
	startPacket.Data = h.Sum(nil)

	return
}

// Bytes marshals start packet and returns byte slice
func (s StartPacket) Bytes() (b []byte) {
	b = append(b, s.Type)
	b = append(b, s.Data...)
	return
}

// Unmarshal unmarshals start packet from byte slice
func (s *StartPacket) Unmarshal(b []byte) (err error) {

	if len(b) <= 1 {
		err = fmt.Errorf("incorrect input length")
		return
	}
	s.Type = b[0]
	s.Data = b[1:]
	return
}
