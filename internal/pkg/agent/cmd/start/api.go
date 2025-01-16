package start

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
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
