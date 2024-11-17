package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bmehdi777/moon/internal/pkg/server/database"
	"gorm.io/gorm"
)

func httpServe(inChannel <-chan *http.Response, outChannel chan<- *http.Request, db *gorm.DB) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleAllRequest(w, r, inChannel, outChannel, db)
	})

	err := http.ListenAndServe(":8080", nil)
	return err
}

func handleAllRequest(w http.ResponseWriter, r *http.Request, inChannel <-chan *http.Response, outChannel chan<- *http.Request, db *gorm.DB) {
	// have to ensure request are going to the right chan
	// so we need to process the domain name
	hostRequest := r.URL.Host
	userInfo := database.FindUserByDomainName(hostRequest, db)
	if userInfo == nil {
		// domain name doesnt exist
		return
	}

	// out and in channel should be in dictionary
	// fqdn represent the id while the value will be channel
	// this dictionary has to be global

	outChannel <- r
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
		fmt.Println("[ERROR]", err)
		return
	}
	r.Body.Close()

	fmt.Println("Body : ", string(bodyBytes))
}
