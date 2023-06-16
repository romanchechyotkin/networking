package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Println("Please provide a port")
		return
	}

	PORT := ":" + args[1]

	listener, err := net.Listen("tcp", PORT)
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The server running on %s port\n", args[1])

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.TrimSpace(string(netData)) == "STOP" {
				fmt.Println("Exiting TCP server!")
				return
		}

		fmt.Printf("From %s\n", conn.RemoteAddr().String())
		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		conn.Write([]byte(myTime))
	}
}