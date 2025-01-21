package utils

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

func HttpHandleAssets(w http.ResponseWriter, r *http.Request, fs fs.FS, pathUrl string) {
	filePath := path.Clean(r.URL.Path)
	if filePath == pathUrl {
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
