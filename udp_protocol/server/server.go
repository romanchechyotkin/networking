package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide a port")
		return
	}

	PORT := ":" + args[1]

	udpAddress, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp4", udpAddress)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The server running on %s port\n", args[1])

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		msg := string(buffer[:n])
		fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), msg)
		
		if strings.TrimSpace(msg) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		data := []byte(reverse(msg))
		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func reverse(msg string) string {
	text := strings.Split(msg, "\n")
	data := []byte(text[0])
	res := make([]byte, len(data))

	for i := len(data)-1; i>=0; i-- {
		res = append(res, data[i])
	}

	return string(res)
}
