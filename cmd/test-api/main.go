package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received req")
		w.Write([]byte("Hello world"))
	})

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		os.Exit(1)
	}
}
