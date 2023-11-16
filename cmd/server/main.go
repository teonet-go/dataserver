// TCP server example
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/teonet-go/teoss/dataserver"
	"github.com/teonet-go/teoss/dataserver/server"
)

const (
	appShort   = "Simple TCP server"
	appVersion = "0.0.1"
	localAddr  = ":8089"
)

func main() {

	// Show application logo
	fmt.Println(appShort + " ver " + appVersion)

	// Set microseconds in log
	log.SetFlags(log.Default().Flags() + log.Lmicroseconds)

	// Create new Data Server
	ds, err := server.NewDataServer(localAddr)
	if err != nil {
		fmt.Println("start data server error:", err)
	}

	for {
		request := dataserver.StartPacket{
			Type: dataserver.READ,
			Data: []byte("start"),
		}
		ds.SetReadRequest(request, func(request dataserver.StartPacket, reader io.Reader, err error) {
			if err != nil {
				fmt.Println("got error:", err)
				return
			}

			p := make([]byte, server.ChankPacketLength)
			for {
				n, err := reader.Read(p)
				if err != nil {
					if err == io.EOF {
						err = nil
					}
					break
				}
				fmt.Printf("got data, n, err: '%s' %d %v\n", string(p[:n]), n, err)
			}
			log.Printf("done '%s', err: %v\n\n", string(request.Data), err)
		})
	}

	select {}

	for {
		// Set Request for read file and get result string in callback
		request := dataserver.StartPacket{
			Type: dataserver.READ,
			Data: []byte("start"),
		}
		ds.SetReadRequest(request, func(buf *bytes.Buffer, err error) {
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Printf("\ngot result:\n%s\n", buf.String())
		})
	}

	// select {}
}
