# Data transmission via TCP

There is golang package to transfer data from client to server
and from server to client via TCP.

[![GoDoc](https://godoc.org/github.com/teonet-go/dataserver?status.svg)](https://godoc.org/github.com/teonet-go/dataserver/)
[![Go Report Card](https://goreportcard.com/badge/github.com/teonet-go/dataserver)](https://goreportcard.com/report/github.com/teonet-go/dataserver)

## Server

Server package implements a TCP server for handling custom data requests and responses.

It allows registering callback functions that will be executed when a connection with a matching start packet is accepted. This allows customizing the server's behavior for different types of requests.

The main inputs are:

1. The local address to listen on
2. The registered callback functions and their associated start packets

The main outputs are:

1. Calling the callback functions when a connection with a matching start packet is accepted.

The code works by:

1. Creating a TCP listener on the provided local address.

2. Launching a goroutine to accept incoming connections.

3. When a connection is accepted, reading a fixed-length start packet.

4. Looking up the callback function registered for that start packet.

5. Calling the callback, passing it the connection and any other needed parameters.

The start packet is used to identify what type of request the connection is for. The callback functions allow custom logic to be executed to handle each request type.

Some key flows are:

- Registering callback functions along with associated start packets on server startup.

- Accepting connections and reading start packets.

- Looking up callback based on start packet and executing it.

- Callback function reading/writing data over the connection.

So in summary, it provides a way to accept connections for custom request types based on start packets, and execute custom logic to handle each type of request. The callbacks allow customizing the input/output and processing for each request.

## Client

Client package implements a TCP client for communicating with a dataserver in Go.

The client.go file defines a Client package that provides functionality for a TCP client to connect to and interact with a dataserver.

It starts by defining some constants:

- startPacketLength - the length in bytes of the start packet sent by client to initialize the connection

- ChunkPacketLength - the max length in bytes of each data packet sent over the connection

- timeout - the read/write timeout for network operations

It then defines a DataClient struct which contains the remote address to connect to, the start packet to send, and the network connection itself.

The NewDataClient function creates a new DataClient instance - it connects to the provided remoteAddr dataserver, sends the startPacket to initialize the connection, and returns the connected DataClient object.

The DataClient has Write and Read methods to send data to and receive data from the dataserver over the TCP connection. Write sends the provided byte slice to the connection. Read reads into the provided byte slice from the connection.

So in summary, the client.go code implements a TCP client that can connect to a dataserver, initialize the connection, and then write data to and read data from the server over the established TCP connection. It provides networking client functionality to interact with a dataserver service.

## Examples

### Server example

[Server example code](cmd/server/main.go) implements a simple TCP server example using the Dataserver package.

It first imports required packages like dataserver, log, fmt, io, etc.

Then it declares some constants like the app name, version, local address to bind the server to.

The main function is the entry point. It prints the app name and version.

It creates a new DataServer instance binding it to the local address.

It then enters an infinite loop to handle client connections continuously.

Each loop iteration, it creates a buffer to hold request data.

It makes a dataserver StartPacket for a READ request with some metadata.

This is passed to the SetRequest method to initialize a reader stream from client.

The callback handles reading data in chunks from the stream into the buffer.

It prints the read data and any errors.

Once reading is done, it makes a new StartPacket for a WRITE request.

This initializes a writer stream to client.

The callback writes the buffered data back to the client in chunks.

It prints the written data and errors.

So in summary, it implements a TCP echo server - clients connect, write data which is read and buffered by the server, and then written back to the client. The dataserver package handles the network communication while the callbacks process the data.

### Client example

[Client example code](cmd/client/main.go) is a Go program that demonstrates a TCP client for a Dataserver.

It starts by defining some constants like the application name, version, and the remote address of the dataserver it will connect to.

The main function is the entry point. It first prints the application name and version. Then it enters a loop that will run once.

Inside the loop, it creates a buffer with some test data to write to the server. It has a hello message repeated multiple times.

It then connects to the dataserver by creating a new DataClient, passing the remote address and a start packet. This initializes the TCP connection.

Next, it writes the buffer data to the dataserver in chunks using the Write method on the DataClient. This sends the data over the TCP connection.

After writing, it closes the client which closes the TCP connection.

Then it re-connects and reads back any data, printing out the chunks it receives. This shows it can send data, disconnect, then connect again and retrieve what was sent.

Finally it closes the client again to clean up the TCP connection.

The purpose is to demonstrate the client API for connecting and communicating with a dataserver over TCP. It shows sending data, then receiving data in a simple example. The key inputs are the test data buffer and remote address. The outputs are the printouts of send/receive status. The logic handles connecting, writing, closing, reconnecting, reading, and closing the TCP client.

## License

[BSD](LICENSE)
