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

	clientManagerWrapper := wrapper.NewLobbyWrapper()

	router := master.NewRouter()
	router.Initialise(clientManagerWrapper)

	router.Start(PORT)
}

func getSettings() (LOG_LEVEL log.Level, PORT int) {
	switch os.Getenv("WARTEMIS_ENV") {
		case "PROD":
			LOG_LEVEL = log.InfoLevel
			PORT = 80
		case "TEST":
			LOG_LEVEL = log.InfoLevel
			PORT = 80
		default:
			// LOG_LEVEL = log.DebugLevel
			LOG_LEVEL = log.InfoLevel
			PORT = 8080
	}
	return
}
