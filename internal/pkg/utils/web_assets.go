package utils

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

func HttpHandleAssets(w http.ResponseWriter, r *http.Request, fs fs.FS, prefixUrl string) {
	filePath := path.Clean(r.URL.Path)
	prefixClean := path.Clean(prefixUrl)

	if filePath == prefixClean {
		filePath = "index.html"
	} else {
		filePath = strings.TrimPrefix(filePath, prefixUrl)
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

	// Change url before giving corresponding file : we don't have prefix in 
	// the FS
	r.URL.Path = filePath

	http.FileServer(http.FS(fs)).ServeHTTP(w, r)
}
