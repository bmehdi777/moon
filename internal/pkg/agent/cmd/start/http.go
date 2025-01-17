package start

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
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
		handleTunnelStatistics(w, r, statistics)
	})

	mux.HandleFunc("/api/healthcheck", healthcheck)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webappHandler(w, r, assets)
	})

	fmt.Println("Dashboard is accessible at : http://localhost:9009")
	err = http.ListenAndServe(":9009", mux)
	if err != nil {
		panic(err)
	}
}
