package main

import (
	"io"
	"log"
	"net"
)


func echo(conn net.Conn){
	defer conn.Close()
	if _, err:= io.Copy(conn, conn); err != nil {
		log.Fatalln("Unable to read/write data")
	}
}

func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Println("Listening...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		log.Println("Accepted connection")

		go echo(conn)
	}
}