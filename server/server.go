// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Server package provides the TCP server implementation and handlers.
package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/teonet-go/dataserver"
)

const (
	// startPacketLength is the length in bytes of the start packet used to
	// initialize new connections.
	startPacketLength = dataserver.StartPacketLength

	// ChankPacketLength is the max length in bytes of data chunks sent after
	// the start packet.
	ChankPacketLength = 21

	// timeout is the read/write timeout for connections.
	timeout = 5 * time.Second
)

var (
	// ErrTimeout is returned when a read or write times out.
	ErrTimeout = errors.New("timeout")

	// ErrIncorrectStartPacket is returned when the start packet is malformed.
	ErrIncorrectStartPacket = errors.New("incorrect start packet")

	// ErrExistsStartPacket is returned when the start packet already exists.
	ErrExistsStartPacket = errors.New("start packet already exists")

	// ErrIncorrectStartPacketLength is returned when the start packet length is incorrect.
	ErrIncorrectStartPacketLength = errors.New("incorrect packet length")
)

// DataServer is the main server instance. It contains the TCP listener,
// request handler map, and sync mutex.
type DataServer struct {
	ln net.Listener
	m  DataServerMap
	*sync.RWMutex
}
type DataServerMap map[DataServerRequest]interface{}
type DataServerRequest [startPacketLength]byte

// NewDataServer creates a new DataServer instance by initializing a DataServer
// struct, creating a new sync.RWMutex, and starting the TCP listener.
// It accepts a localAddr string to listen on and returns a DataServer pointer
// and error.
func NewDataServer(localAddr string) (ds *DataServer, err error) {
	ds = &DataServer{m: make(DataServerMap), RWMutex: new(sync.RWMutex)}
	ds.ln, err = ds.listening(localAddr)
	return
}

// LocalPort returns the local port number that the DataServer is listening on.
func (ds DataServer) LocalPort() (port int) {
	if addr, ok := ds.ln.Addr().(*net.TCPAddr); ok {
		port = addr.Port
	}
	return
}

// SetRequest registers read or write request depending on the type of start
// packet. Callback function calls when connection accepted.
//
// Avalable callback functions:
//
//	// Read callback
//	func(startPacket *dataserver.StartPacket, reader io.Reader, err error)
//
//	// Write callback
//	func(startPacket *dataserver.StartPacket, writer io.Writer, err error)
//
//	// ReadAll calback
//	func(buf *bytes.Buffer, err error)
func (ds *DataServer) SetRequest(startPacket *dataserver.StartPacket,
	callback interface{}) (err error) {
	b := startPacket.Bytes()

	// Check start packet length
	if len(b) != startPacketLength {
		err = ErrIncorrectStartPacketLength
		return
	}

	// Check start packet already registered
	if ok := ds.check(b); ok {
		err = ErrExistsStartPacket
		return
	}

	// Register start packet and return
	ds.add(b, callback)
	return
}

// listening starts a TCP listener for the server on the provided local address
// and port. It launches a goroutine to accept incoming connections and handle
// them by calling ds.handleConnection.
func (ds DataServer) listening(localAddr string) (ln net.Listener, err error) {

	// Start listening
	ln, err = net.Listen("tcp", localAddr)
	if err != nil {
		return
	}

	// Accept incoming connections and handle them
	go func() {
		for {
			// Accept incoming connections
			conn, err := ln.Accept()
			if err != nil {
				log.Println("incomming connection error:", err)
				continue
			}

			// Handle the connection in a new goroutine
			go ds.handleConnection(conn)
		}
	}()

	return
}

// handleConnection handles an incoming client connection.
// It reads the start packet, validates it, looks up the
// registered callback function, and executes the callback
// based on its type.
func (ds DataServer) handleConnection(conn net.Conn) {

	// Close the connection when we're done
	defer conn.Close()

	var buf = new(bytes.Buffer)
	var err error

	// Read start packet
	waitStartPacket := make(chan error, 1)
	startPacketBuf := make([]byte, startPacketLength)
	go func() {
		_, err := conn.Read(startPacketBuf)
		if err == nil {
			// Check received start packet
			if !ds.check(startPacketBuf) {
				err = ErrIncorrectStartPacket
			}
		}
		waitStartPacket <- err
	}()

	// Answer on return
	defer func() { ds.del(startPacketBuf) }()

	// Wait start packet during timeout
	select {
	case err = <-waitStartPacket:
	case <-time.After(timeout):
		err = ErrTimeout
	}
	if err != nil {
		log.Println("got start packet error:", err)
		return
	}

	// Unmarshal start packet buffer
	startPacket := &dataserver.StartPacket{}
	err = startPacket.Unmarshal(startPacketBuf)
	if err != nil {
		err = fmt.Errorf("start packet unmarshal error: %s", err)
		log.Println(err)
		return
	}

	// Get callback function registerred for this start packet
	callback, ok := ds.get(startPacketBuf)
	if !ok {
		callback = nil
	}

	// Execute callback
	switch cb := callback.(type) {

	// Read All incoming data to buffer callback
	case func(buf *bytes.Buffer, err error):
		total := 0
		b := make([]byte, ChankPacketLength)
		for {
			n, err := conn.Read(b)
			if err != nil {
				if err == io.EOF {
					// Done successfully
					err = nil
				}
				break
			}
			if callback != nil {
				buf.Write(b[:n])
			}
			total += n
		}
		cb(buf, err)

	// Execute accepted read callback
	case func(startPacket *dataserver.StartPacket, reader io.Reader, err error):
		cb(startPacket, conn, err)

	// Execute accepted write callback
	case func(startPacket *dataserver.StartPacket, writer io.Writer, err error):
		cb(startPacket, conn, err)

	default:
		log.Fatalln("some other callback", cb)
	}
}

// add adds request to DataServerMap
func (ds *DataServer) add(request []byte, res interface{}) {
	ds.Lock()
	defer ds.Unlock()
	ds.m[DataServerRequest(request)] = res
}

// del deletes request from DataServerMap
func (ds *DataServer) del(request []byte) {
	ds.Lock()
	defer ds.Unlock()
	delete(ds.m, DataServerRequest(request))
}

// check returns true if request exists in map
func (ds *DataServer) check(request []byte) (ok bool) {
	ds.RLock()
	defer ds.RUnlock()
	_, ok = ds.m[DataServerRequest(request)]
	return
}

// get returns the callback function and ok bool for the given request if it
// exists in the request map ds.m. It locks the read lock, defer unlocks, and
// looks up the request in ds.m to find the callback.
func (ds *DataServer) get(request []byte) (callback interface{}, ok bool) {
	ds.RLock()
	defer ds.RUnlock()
	callback, ok = ds.m[DataServerRequest(request)]
	return
}
