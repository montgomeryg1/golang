package main

import (
	"io"
	"log"
	"net"
)


func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "www.trurating.com:443")
	if err != nil {

	}
	defer dst.Close()
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":443")
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

		go handle(conn)
	}
}