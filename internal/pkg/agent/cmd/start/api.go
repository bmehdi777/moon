package start

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func webappHandler(w http.ResponseWriter, r *http.Request, fs fs.FS) {
	filePath := path.Clean(r.URL.Path)
	if filePath == "/" {
		filePath = "index.html"
	} else {
		filePath = strings.TrimPrefix(filePath, "/")
	}

	file, err := fs.Open(filePath)
	if os.IsNotExist(err) || filePath == "index.html" {
		http.ServeFileFS(w, r, fs, "index.html")
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	http.FileServer(http.FS(fs)).ServeHTTP(w, r)
}

func handleTunnelStatistics(w http.ResponseWriter, r *http.Request, statistics *Statistics) {
	// SSE config
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientGone := r.Context().Done()

	fmt.Println("Connection open")

	fmt.Fprintf(w, "data: %s\n\n", statistics.HttpCalls.Format())
	w.(http.Flusher).Flush()

	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <-clientGone:
			fmt.Println("Connection closed")
			return
		case <-tick:
			w.(http.Flusher).Flush()
		case <-statistics.Event:
			fmt.Fprintf(w, "data: %s\n\n", statistics.HttpCalls.Format())
			w.(http.Flusher).Flush()
		}
	}
}
