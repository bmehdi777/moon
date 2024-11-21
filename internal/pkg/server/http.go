package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"github.com/bmehdi777/moon/internal/pkg/server/database"
	"gorm.io/gorm"
)

func httpServe(channelsPerDomain ChannelsDomains, db *gorm.DB) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleAllRequest(w, r, channelsPerDomain, db)
	})

	fullAddrFmt := fmt.Sprintf("%v:%v", config.GlobalConfig.TcpAddr, config.GlobalConfig.HttpPort)
	log.Printf("HTTP server is up at %v", fullAddrFmt)
	err := http.ListenAndServe(fullAddrFmt, nil)
	return err
}

func handleAllRequest(w http.ResponseWriter, r *http.Request, channelsPerDomain ChannelsDomains, db *gorm.DB) {
	// have to ensure request are going to the right chan
	// so we need to process the domain name
	hostRequest := r.URL.Host
	userInfo := database.FindUserByDomainName(hostRequest, db)
	if userInfo == nil {
		// domain name doesnt exist
		return
	}

	// out and in channel should be in dictionary
	// fqdn represent the id while the value will be channel
	// this dictionary has to be global

	// it's working - just not now

	//outChannel <- r
	//response := <-inChannel
	//w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	//w.Header().Set("Content-Length", response.Header.Get("Content-Length"))
	//io.Copy(w, response.Body)
	//response.Body.Close()
}
