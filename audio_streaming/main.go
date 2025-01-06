package main

import (
    "log"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./data"))

    http.Handle("/dash/", http.StripPrefix("/dash/", fs))

    log.Println("Starting server on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
