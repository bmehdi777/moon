package start

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"
)

//go:embed all:dist
var distFolder embed.FS

func handleHttpServer(statistics *Statistics) {
	assets, err := fs.Sub(distFolder, "dist")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/tunnels/status", func(w http.ResponseWriter, r *http.Request) {
		tunnelStatus(w, r, statistics)
	})

	mux.HandleFunc("/api/healthcheck", healthcheck)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webappHandler(w, r, assets)
	})

	err = http.ListenAndServe(":9009", mux)
	if err != nil {
		panic(err)
	}
}

func tunnelStatus(w http.ResponseWriter, r *http.Request, statistics *Statistics) {
	// SSE config
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientGone := r.Context().Done()

	fmt.Println("Connection open")

	fmt.Fprintf(w, "data: hello\n\n")
	w.(http.Flusher).Flush()

	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <-clientGone:
			fmt.Println("Connection closed")
			return
		case <-tick:
			fmt.Println("Heartbeat")
			w.(http.Flusher).Flush()
		case <-statistics.Event:
			fmt.Println("Req : ", statistics.HttpCalls[0].Request.URL)
			fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf(`{"req_uri": "test"}`))
			w.(http.Flusher).Flush()
		}
	}
}

