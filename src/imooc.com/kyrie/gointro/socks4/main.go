package main

import (
	"net"
	"log"
	"io"
)

func main() {
	listener, err := net.Listen("tcp", ":2000")
	panicOnError(err)
	defer listener.Close()
	for {
		conn, err := listener.Accept();
		if err != nil {
			continue
		}
		go relay(conn)
	}
}
func relay(connection net.Conn) {
	defer connection.Close()
	log.Printf("relay: listen on %v, new client from %v\n", connection.LocalAddr().String(), connection.RemoteAddr().String())
	serverConnection, err := net.Dial("tcp", "localhost:3000")
	panicOnError(err)
	defer serverConnection.Close()
	log.Printf("relay: connected to %v, from %v\n", serverConnection.LocalAddr().String(), serverConnection.RemoteAddr().String())
	done := make(chan bool)
	go func() {
		forward(connection, serverConnection)
		done<-true
	}()
	forward(serverConnection, connection)
	<-done
}
func forward(destination net.Conn, source net.Conn) {
	io.Copy(destination, source)
	log.Printf("relay: done copying from %v to %v\n", source.LocalAddr().String(), destination.RemoteAddr().String())
	tcpConnection := destination.(*net.TCPConn)
	tcpConnection.CloseWrite()
}

func panicOnError(err error) {
	if err !=nil {
		panic(err);
	}
}

