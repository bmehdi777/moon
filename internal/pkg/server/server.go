package server

import (
	"log"
	"net/http"
	"os"

	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"github.com/bmehdi777/moon/internal/pkg/server/database"
)

func Run() {
	config.InitConfig()

	db, err := database.InitializeDBConn()
	if err != nil {
		log.Fatalf("Can't connect to the database.")
		os.Exit(1)
	}

	channelPerDomain := make(ChannelsDomains)

	// tcp connection between client and server
	go tcpServe(&channelPerDomain, db)

	// http connection between server and the world
	if err := httpServe(&channelPerDomain, db); err != nil {
		log.Fatalf("Error : %v ", err)
	}
}

type ChannelsHttp struct {
	RequestChannel  chan *http.Request
	ResponseChannel chan *http.Response
}

type ChannelsDomains map[string]ChannelsHttp

func (c *ChannelsDomains) Add(name string) {
	(*c)[name] = ChannelsHttp{
		RequestChannel:  make(chan *http.Request),
		ResponseChannel: make(chan *http.Response),
	}
}

func (c *ChannelsDomains) Delete(name string) {
	delete(*c, name)
}

func (c *ChannelsDomains) Get(name string) *ChannelsHttp {
	if channel, found := (*c)[name]; found {
		return &channel
	}
	return nil
}
