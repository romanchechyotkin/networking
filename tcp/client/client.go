package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}

	CONNECT := args[1]

	conn, err := net.Dial("tcp", CONNECT)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The TCP server is %s\n", conn.RemoteAddr().String())

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		fmt.Fprintf(conn, string(data))

		msg, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("-> ", msg)
		
		if strings.TrimSpace(string(data)) == "STOP" {
			log.Println("TCP client exiting...")
			return
		}
	}
}