package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// testing purpose
		if r.Method == "POST" {
			time.Sleep(5*time.Second)
		}
		fmt.Println("Received req")
		w.Write([]byte("Hello world"))
	})

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		os.Exit(1)
	}
}
