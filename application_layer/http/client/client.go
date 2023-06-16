package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	
	if len(args) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}

	CONNECT := "http://" + args[1]

	resp, err := http.Get(CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	
	fmt.Println(body)
	fmt.Println(string(body))
}