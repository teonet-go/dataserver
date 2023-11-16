package dataserver

import "fmt"

const (
	READ = iota
	WRITE
)

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

// Bytes marshals start packet and returns byte slice
func (s StartPacket) Bytes() (b []byte) {
	b = append(b, s.Type)
	b = append(b, s.Data...)
	return
}

// Unmarshal unmarshals start packet from byte slice
func (s *StartPacket) Unmarshal(b []byte) (err error) {

	if len(b) <= 1 {
		err = fmt.Errorf("wrong number of bytes")
		return
	}
	s.Type = b[0]
	s.Data = b[1:]
	return
}
