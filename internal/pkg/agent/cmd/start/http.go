package start

import (
	"encoding/json"
	"io"
	"net/http"
)

func httpServe(stats *Statistics) error {
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		handleStats(w, r, stats)
	})

	err := http.ListenAndServe("localhost:5555", nil)
	return err
}

func handleStats(w http.ResponseWriter, _ *http.Request, stats *Statistics) {
	statsStrings, err := json.Marshal(stats)
	if err != nil {
		io.WriteString(w, err.Error())
	}

	w.Write(statsStrings)
}
