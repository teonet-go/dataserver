// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dataserver TCP server example.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/teonet-go/dataserver"
	"github.com/teonet-go/dataserver/server"
)

const (
	appShort   = "Simple TCP server" // Application short name.
	appVersion = "0.0.1"             // Application version.
	localAddr  = ":8089"             // Local address of the dataserver.
)

// main is the entry point for the Dataserver TCP Server example. It initializes
// the dataserver, sets up read and write requests in a loop with callbacks,
// and handles errors.
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
		var buf = new(bytes.Buffer)

		// Set read request
		startPacket := dataserver.MakeStartPacket(dataserver.READ, "start", -1)
		err := ds.SetRequest(startPacket,
			func(startPacket *dataserver.StartPacket, reader io.Reader, err error) {

				if err != nil {
					fmt.Println("an error occurred while initializing the reader:", err)
					return
				}

				fmt.Println("ready to read data")

				p := make([]byte, server.ChankPacketLength)
				for {
					n, err := reader.Read(p)
					if err != nil {
						if err == io.EOF {
							err = nil
						}
						break
					}
					buf.Write(p[:n])
					fmt.Printf("got data, n, err: '%s' %d %v\n",
						string(p[:n]), n, err)
				}
				log.Printf("done %x, err: %v\n\n", startPacket.Data, err)

				// Set write request
				startPacket = dataserver.MakeStartPacket(dataserver.WRITE, "start", -1)
				ds.SetRequest(startPacket,
					func(startPacket *dataserver.StartPacket, writer io.Writer, err error) {

						if err != nil {
							fmt.Println("an error occurred while initializing the writer:", err)
							return
						}

						fmt.Println("ready to write data")
						for {
							n, err := buf.Read(p)
							if err != nil {
								break
							}
							writer.Write(p[:n])
							fmt.Printf("send data, n, err: '%s' %d %v\n",
								string(p[:n]), n, err)
						}
						log.Printf("done %x, err: %v\n\n", startPacket.Data, err)
					},
				)

			},
		)
		if err != nil {
			if err == server.ErrExistsStartPacket {
				time.Sleep(250 * time.Millisecond)
				continue
			}
			log.Println("set read request error:", err)
			return
		}
	}
}
