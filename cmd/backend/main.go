package main

import (
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/base"
	"github.com/Project-Wartemis/pw-backend/internal/master"
	"github.com/Project-Wartemis/pw-backend/internal/http"
)

func main() {
	LOG_LEVEL, PORT := getSettings()

	log.SetLevel(LOG_LEVEL)

	log.Info("Execute main")

	lobby := base.NewLobby()

	lobbyHttpInterface := http.NewLobbyHttpInterface(lobby)

	router := master.NewRouter()
	router.Initialise(lobbyHttpInterface)

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
