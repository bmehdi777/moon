package server

import (
	"log"
	"os"

	"moon/internal/pkg/server/config"
	"moon/internal/pkg/server/database"
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
