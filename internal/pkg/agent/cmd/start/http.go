package start

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed all:dist
var distFolder embed.FS

func handleHttpServer() {
	assets, err := fs.Sub(distFolder, "dist")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/", apiHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webappHandler(w, r, assets)

	})

	err = http.ListenAndServe(":9009", mux)
	if err != nil {
		panic(err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
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
