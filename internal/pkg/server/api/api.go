package api

import (
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	HealthcheckHandler  Healthcheck
	CertificatesHandler Certificates
}

func NewApp() App {
	return App{
		HealthcheckHandler:  Healthcheck{},
		CertificatesHandler: Certificates{},
	}
}

func (a *App) ServeHttp(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/api", "", 1)

	var certsId string

	switch {
	case match(path, "/health"):
		a.HealthcheckHandler.handleGet(w)
	case match(path, "/certificates") && r.Method == "POST":
		a.CertificatesHandler.handleCreate(w)
	case match(path, "/certificates/+", &certsId) && r.Method == "GET":
		a.CertificatesHandler.handleGet(certsId, w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func match(path, pattern string, vars ...interface{}) bool {
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '+':
			// '+' matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			switch p := vars[0].(type) {
			case *string:
				*p = segment
			case *int:
				n, err := strconv.Atoi(segment)
				if err != nil || n < 0 {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
}
