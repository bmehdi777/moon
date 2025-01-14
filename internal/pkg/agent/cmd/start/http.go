package start

import (
	"embed"
	"io/fs"
	"net/http"
)

// Files aren't here before compiling

//go:embed all:dist
var distFolder embed.FS

func handleHttpServer() {
	assets, err := fs.Sub(distFolder, "dist")
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.FS(assets))
	http.Handle("/", http.StripPrefix("/", fs))

	err = http.ListenAndServe(":9009", nil)
	if err != nil {
		panic(err)
	}
}
