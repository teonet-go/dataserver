package dataserver

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

const (
	READ  = iota // Receive data on data server request.
	WRITE        // Send data from data server request.
)

const StartPacketLength = 25

// StartPacket contains start packet data and type, and methods to murshal &
// unmarshal it.
type StartPacket struct {
	// Type of start packet:
	//  0 - READ (server will receive data);
	//  1 - WRITE (server will send data).
	Type byte

	// Object size
	ObjectSize int64

	// Start packet data
	Data []byte
}

// MakeStartPacket creates new start packet from Type and Name
func MakeStartPacket(t byte, name string, objectSize int64) (startPacket *StartPacket) {

	startPacket = &StartPacket{Type: t, ObjectSize: objectSize}

	h := md5.New()
	h.Write([]byte(name))
	startPacket.Data = h.Sum(nil)

	return
}

// Bytes marshals start packet and returns byte slice
func (s StartPacket) Bytes() (b []byte) {

	b = append(b, s.Type)

	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(s.ObjectSize))
	b = append(b, size...)

	b = append(b, s.Data...)
	return
}

// Unmarshal unmarshals start packet from byte slice
func (s *StartPacket) Unmarshal(b []byte) (err error) {

	if len(b) <= 9 {
		err = fmt.Errorf("incorrect input length")
		return
	}
	s.Type = b[0]

	s.ObjectSize = int64(binary.LittleEndian.Uint64(b[1:9]))

	s.Data = b[9:]
	return
}
