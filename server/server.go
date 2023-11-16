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
	startPacketLength = 6
	ChankPacketLength = 21
	timeout           = 5 * time.Second
)

var (
	ErrTimeout              = errors.New("timeout")
	ErrIncorrectStartPacket = errors.New("incorrect start packet")
)

// DataServer is TCP Data Server data structure and methods receiver.
type DataServer struct {
	ln net.Listener
	m  DataServerMap
	*sync.RWMutex
}
type DataServerMap map[DataServerRequest]interface{}
type DataServerRequest [startPacketLength]byte

type DataServerReadAllCallback func(buf *bytes.Buffer, err error)
type DataServerAcceptedCallback func(request dataserver.StartPacket, reader io.Reader, err error)

type DataServerReader struct {
	net.Conn
}

// NewDataServer creates a new DataServer object.
func NewDataServer(localAddr string) (ds *DataServer, err error) {
	ds = &DataServer{m: make(DataServerMap), RWMutex: new(sync.RWMutex)}
	ds.ln, err = ds.listening(localAddr)
	return
}

// SetReadRequest sets read and returns reader to get data
func (ds *DataServer) SetReadRequest(request dataserver.StartPacket,
	res interface{}) {
	ds.add(request.Bytes(), res)
}

// listening starts listening and accept incoming connections.
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
			log.Printf("connected from: %s\n", conn.RemoteAddr())

			// Handle the connection in a new goroutine
			go ds.handleConnection(&DataServerReader{conn})
		}
	}()

	return
}

// handleConnection handles the client connection.
func (ds DataServer) handleConnection(dsr *DataServerReader) {

	// Close the connection when we're done
	defer dsr.Conn.Close()

	var buf = new(bytes.Buffer)
	var err error

	// Read start packet
	waitStartPacket := make(chan error, 1)
	startPacketBuf := make([]byte, startPacketLength)
	go func() {
		_, err := dsr.Conn.Read(startPacketBuf)
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
	log.Println("start packet received")

	callback, ok := ds.get(startPacketBuf)
	if !ok {
		callback = nil
	}

	switch cb := callback.(type) {

	// Read All incoming data by chanks to buffer
	// DataServerReadAllCallback
	case func(buf *bytes.Buffer, err error):
		total := 0
		b := make([]byte, ChankPacketLength)
		for {
			n, err := dsr.Conn.Read(b)
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

	// Execute accepted callback
	// DataServerAcceptedCallback
	case func(request dataserver.StartPacket, reader io.Reader, err error):
		startPacket := dataserver.StartPacket{}
		err := startPacket.Unmarshal(startPacketBuf)
		if err != nil {
			err = fmt.Errorf("startPacket.Unmarshal error: %s", err)
		}
		cb(startPacket, dsr, err)

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

// get returns requests callback function
func (ds *DataServer) get(request []byte) (callback interface{}, ok bool) {
	ds.RLock()
	defer ds.RUnlock()
	callback, ok = ds.m[DataServerRequest(request)]
	return
}
