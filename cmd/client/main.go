// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dataserver TCP client example.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/teonet-go/dataserver"
	"github.com/teonet-go/dataserver/client"
)

const (
	appShort   = "Simple TCP client" // Application short name.
	appVersion = "0.0.1"             // Application version.
	remoteAddr = "localhost:8089"    // Remote address of the dataserver.
)

// main is the entry point for the Dataserver TCP Client example. It connects
// to a dataserver, writes some data, then reads back data. It demonstrates
// using the dataserver client API.
func main() {

	// Show application logo
	fmt.Println(appShort + " ver " + appVersion)

	for j := 1; j <= 1; j++ {

		j = 4

		// Create data buffer
		var buf = new(bytes.Buffer)
		for i := 0; i < 10; i++ {
			buf.Write([]byte(fmt.Sprintf("Hello, server - %d-%d!", j, i)))
		}
		buf.Write([]byte("Bye!"))

		// Title
		title := func(dc *client.DataClient, t string) (err error) {
			fmt.Println(t)
			_, err = io.WriteString(dc, t)
			return
		}

		// Write data to server.
		// Connect to Data Server and send start packet (request).
		startPacket := dataserver.MakeStartPacket(dataserver.READ, "start", -1)
		dc, err := client.NewDataClient(remoteAddr, startPacket)
		if err != nil {
			log.Println("can't connect to data server, error: ", err)
			return
		}
		fmt.Printf("\nconnected to: %s\n", dc.RemoteAddr())

		// Execute write data by case
		switch j {

		// Write using read write by chanks
		case 4:
			title(dc, "write using read write by chanks")
			p := make([]byte, client.ChankPacketLength)
			var n int
			for {
				// Get data from buffer
				n, err = buf.Read(p)
				if err != nil {
					if err == io.EOF {
						err = nil
					}
					break
				}

				// Write data to DataClient connection
				_, err = dc.Write(p[:n])
				if err != nil {
					break
				}
			}
			dc.Close()
		}

		// Read data from server.
		// Connect to Data Server and send start packet (request).
		startPacket = dataserver.MakeStartPacket(dataserver.WRITE, "start", -1)
		dc, err = client.NewDataClient(remoteAddr, startPacket)
		if err != nil {
			log.Println("can't connect to data server, error: ", err)
			return
		}
		fmt.Printf("\nconnected to: %s\n", dc.RemoteAddr())

		// TODO: get data from data server
		fmt.Println("read using read write by chanks")
		p := make([]byte, client.ChankPacketLength)
		var n int
		for {
			// Get data from  DataClient connection
			n, err = dc.Read(p)
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}

			// Print received data
			fmt.Printf("got data, n, err: '%s' %d %v\n",
				string(p[:n]), n, err)
		}

		// CLose connection
		dc.Close()
	}
}
