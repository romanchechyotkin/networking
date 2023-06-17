package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	CONNECT := "pop.gmail.com:995"

	conn, err := net.Dial("tcp", CONNECT)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("The POP server is %s\n", conn.RemoteAddr().String())

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text)
		fmt.Fprintf(conn, string(data))

		msg, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("-> ", msg)
	}
}
