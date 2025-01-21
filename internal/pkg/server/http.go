package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"moon/internal/pkg/server/api"
	"moon/internal/pkg/server/config"
	"moon/internal/pkg/server/database"
	"moon/internal/pkg/utils"

	"gorm.io/gorm"
)

func httpServe(channelsPerDomain *ChannelsDomains, db *gorm.DB) error {

	mux := http.NewServeMux()

	appRouter := api.NewApp()

	tunHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleTunnelRequest(w, r, channelsPerDomain, db)
	})


	assets := os.DirFS(config.GlobalConfig.App.AssetsFolderPath)
	webAssetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.HttpHandleAssets(w, r, assets, "/web/")
	})

	mux.Handle("/api/", middlewareTun(appRouter.ServeHttp, tunHandler))
	// not sure I will keep the "/web/" : should the application be rendered
	// as / or /web/ ? will see i guess
	mux.Handle("/web/", middlewareTun(webAssetsHandler, tunHandler))


	// TODO: may be changed I think : we would want to redirect everything to
	// /web/ 
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
