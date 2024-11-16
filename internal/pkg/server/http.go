package server

import (
	"fmt"
	"io"
	"net/http"
)

func httpServe(inChannel <-chan *http.Response, outChannel chan<- *http.Request) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleAllRequest(w, r, inChannel, outChannel)
	})

	err := http.ListenAndServe(":8080", nil)
	return err
}

func handleAllRequest(w http.ResponseWriter, r *http.Request, inChannel <-chan *http.Response, outChannel chan<- *http.Request) {
	outChannel <- r
	//logRequest(r)
	response := <-inChannel
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", response.Header.Get("Content-Length"))
	io.Copy(w, response.Body)
	response.Body.Close()
}

func logRequest(r *http.Request) {
	fmt.Println("Headers : ")
	for key, value := range r.Header {
		fmt.Println(key, value)
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("[ERROR]", err.Error)
		return
	}
	r.Body.Close()

	fmt.Println("Body : ", string(bodyBytes))
}
