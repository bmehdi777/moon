package server

import (
	"os"

	"moon/internal/pkg/server/config"
	"moon/internal/pkg/server/database"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run() {
	config.InitConfig()

	configureLogger()

	db, err := database.InitializeDBConn()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Can't connect to the database.")
		os.Exit(1)
	}

	channelPerDomain := make(ChannelsDomains)

	// tcp connection between client and server
	go tcpServe(&channelPerDomain, db)

	// http connection between server and the world
	if err := httpServe(&channelPerDomain, db); err != nil {
		log.Fatal().Stack().Err(err)
	}
}

func configureLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
