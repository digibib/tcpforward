package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	localAddr  = flag.String("l", ":9999", "host:port to listen on")
	remoteAddr = flag.String("r", ":9200", "host:port to forward to")
	prefix     = flag.String("p", "tcpforward: ", "String to prefix log output")
)

func forward(conn net.Conn) {
	client, err := net.Dial("tcp", *remoteAddr)
	if err != nil {
		log.Printf("Dial failed: %v", err)
		return
	}
	log.Printf("Forwarding from %v to %v\n", conn.LocalAddr(), client.RemoteAddr())
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func main() {
	flag.Parse()
	log.SetPrefix(*prefix + ": ")

	listener, err := net.Listen("tcp", *localAddr)
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr().String())
		go forward(conn)
	}
}
