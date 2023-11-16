// TCP client example

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
	appShort   = "Simple TCP client"
	appVersion = "0.0.1"
	remoteAddr = "localhost:8089"
)

func main() {

	// Show application logo
	fmt.Println(appShort + " ver " + appVersion)

	// sendData(remoteAddr)
	// select {}

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

		// Connect to Data Server and send start packet (request)
		request := dataserver.StartPacket{
			Type: dataserver.READ,
			Data: []byte("start"),
		}
		dc, err := client.NewDataClientWriter(remoteAddr, request)
		if err != nil {
			log.Println("can't connect to data server, error: ", err)
			return
		}
		fmt.Printf("\nconnected to: %s\n", dc.RemoteAddr())

		// Execute write data by case
		switch j {

		// // Write all data buffer
		// case 1:
		// 	title(dc, "Write all data buffer.")
		// 	err = dc.WriteAll()
		// 	if err != nil {
		// 		log.Println("write all, error:", err)
		// 		return
		// 	}

		// // Write all data as one string
		// case 2:
		// 	title(dc, "Write all data as one string.")
		// 	io.WriteString(dc, buf.String())
		// 	dc.Close()

		// // Write using io.Copy
		// case 3:
		// 	title(dc, "Write using io.Copy.")
		// 	io.Copy(dc.Conn, buf)
		// 	dc.Close()

		// Write using read write by chanks
		case 4:
			title(dc, "Write using read write by chanks.")
			p := make([]byte, client.ChankPacketLength)
			var n int
			for {
				// Get data from reader
				n, err = buf.Read(p) // dc.reader.Read(p)
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
	}

	// select {}
}

// func sendData(remoteAddr string) {

// 	// Connect to the server
// 	conn, err := net.Dial("tcp", remoteAddr)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("connected to: %s\n", conn.RemoteAddr())

// 	// Send start packet to the server
// 	_, err = conn.Write([]byte("start"))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Send some data to the server
// 	for i := 0; i < 10; i++ {
// 		_, err = conn.Write([]byte("Hello, server!\n"))
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}

// 	// Send last data packet which may be less than 14 bytes
// 	_, err = conn.Write([]byte("Bye!\n"))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Close the connection
// 	conn.Close()
// }
