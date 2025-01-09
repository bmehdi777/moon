package server

import (
	"log"
	"os"

	"context"
	"fmt"

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

	go func() {
		ctx := context.Background()
		srv := TunnelServer{
			ChannelsPerDomain: &channelPerDomain,
		}
		err = srv.Run(ctx)
		if err != nil {
			fmt.Println("Err : ", err)
		}
	}()
	// tcp connection between client and server
	//go tcpServe(&channelPerDomain, db)

	// http connection between server and the world
	if err := httpServe(&channelPerDomain, db); err != nil {
		log.Fatalf("Error : %v ", err)
	}

}
