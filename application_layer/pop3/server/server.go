package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type POP3Server struct {
	addr string
}

func main() {
	ADDR := "127.0.0.1:110"

	pop3Server := &POP3Server{
		addr: ADDR,
	}

	log.Printf("POP3 server running on %s\n", pop3Server.addr)
	log.Fatal(pop3Server.ListenAndServe())
}

func (srv *POP3Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go srv.handleConnection(conn)
	}
}

func (srv *POP3Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("Accepted connection from", conn.RemoteAddr().String())
	
	scanner := bufio.NewScanner(conn)

	srv.writeResponse(conn, "+OK POP3 server ready")

	for scanner.Scan() {
		command := scanner.Text()
		log.Println("Received command", command)
	}
}

func (srv *POP3Server) writeResponse(conn net.Conn, text string)  {
	fmt.Fprintf(conn, "%s\r\n", text)
	log.Println("Sent response:", text)
}