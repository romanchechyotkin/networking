package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}

	CONNECT := args[1]

	udpAddress, err := net.ResolveUDPAddr("udp4", CONNECT)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp4", nil, udpAddress)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The UDP server is %s\n", conn.RemoteAddr().String())

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text)
		conn.Write(data)

		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d bytes from %s\n", n, addr.String())
		msg := string(buffer[:n])
		fmt.Printf("%s\n", msg)
	}
}

