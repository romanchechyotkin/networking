package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Println("Please provide a port")
		return
	}

	PORT := ":" + args[1]

	http.HandleFunc("/", index)
	http.ListenAndServe(PORT, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}
