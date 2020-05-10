package main

import (
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/master"
	"github.com/Project-Wartemis/pw-backend/internal/wrapper"
)

func main() {
	LOG_LEVEL, PORT := getSettings()

	log.SetLevel(LOG_LEVEL)

	log.Info("Execute main")

	lobbyWrapper := wrapper.NewLobbyWrapper()
	roomWrapper := wrapper.NewRoomWrapper()

	router := master.NewRouter()
	router.Initialise(lobbyWrapper, roomWrapper)

	router.Start(PORT)
}

func getSettings() (LOG_LEVEL log.Level, PORT int) {
	switch os.Getenv("WARTEMIS_ENV") {
		case "BUILD":
			LOG_LEVEL = log.InfoLevel
			PORT = 80
		default:
			// LOG_LEVEL = log.DebugLevel
			LOG_LEVEL = log.InfoLevel
			PORT = 8080
	}
	return
}
