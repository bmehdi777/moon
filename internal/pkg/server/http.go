package server

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"

	"moon/internal/pkg/server/api"
	"moon/internal/pkg/server/config"
	"moon/internal/pkg/server/database"

	"gorm.io/gorm"
)

//go:embed all:dist
var distFolder embed.FS

func httpServe(channelsPerDomain *ChannelsDomains, db *gorm.DB) error {
	mux := http.NewServeMux()

	appRouter := api.NewApp()

	tunHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleTunnelRequest(w, r, channelsPerDomain, db)
	})

	mux.Handle("/api/", middlewareTun(appRouter.ServeHttp, tunHandler))
	mux.Handle("/web/", middlewareTun(appRouter.ServeHttp, tunHandler))

	mux.HandleFunc("/", tunHandler)

	fullAddrFmt := fmt.Sprintf("%v:%v", config.GlobalConfig.App.HttpAddr, config.GlobalConfig.App.HttpPort)
	log.Printf("HTTP server is up at %v", fullAddrFmt)
	err := http.ListenAndServe(fullAddrFmt, mux)
	return err
}

func middlewareTun(api http.HandlerFunc, tun http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == config.GlobalConfig.App.GlobalDomainName {
			api.ServeHTTP(w, r)
		} else {
			tun.ServeHTTP(w, r)
		}
	})
}

func handleTunnelRequest(w http.ResponseWriter, r *http.Request, channelsPerDomain *ChannelsDomains, db *gorm.DB) {
	// have to ensure request are going to the right chan
	// so we need to process the domain name
	hostRequest := r.Host
	log.Printf("Host url : %v", hostRequest)
	record := database.FindDomainRecordByName(hostRequest, db)
	if record == nil {
		http.Error(w, "Record not found.", http.StatusNotFound)
		return
	}

	if !record.ConnectionOpen {
		http.Error(w, "Record not found.", http.StatusNotFound)
		return
	}

	channel, ok := (*channelsPerDomain)[hostRequest]
	if !ok {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	channel.RequestChannel <- r
	response := <-channel.ResponseChannel
	io.Copy(w, response.Body)
	response.Body.Close()
}
