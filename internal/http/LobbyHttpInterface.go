package http

import (
	"net/http"
	"github.com/gorilla/websocket"
	sync "github.com/sasha-s/go-deadlock"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/base"
)

type LobbyHttpInterface struct {
	sync.RWMutex
	lobby *base.Lobby
	upgrader *websocket.Upgrader
}

func NewLobbyHttpInterface(lobby *base.Lobby) *LobbyHttpInterface {
	upgrader := &websocket.Upgrader {
		CheckOrigin: func(r *http.Request) bool {
			return true // accept connections from anywhere
		},
	}
	return &LobbyHttpInterface {
		lobby: lobby,
		upgrader: upgrader,
	}
}

func (this *LobbyHttpInterface) HandleNewConnection(writer http.ResponseWriter, request *http.Request) {
	conn, err := this.getUpgrader().Upgrade(writer, request, nil)
	if err != nil {
		log.Errorf("Cannot upgrade to websocket: %s", err)
		return
	}
	defer conn.Close()

	connection := base.NewConnection(conn)
	this.getLobby().HandleConnect(connection)
	defer connection.HandleDisconnect()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Unable to read message: %s", err)
			return
		}
		connection.HandleMessage(message)
	}
}



// getters and setters

func (this *LobbyHttpInterface) getLobby() *base.Lobby {
	this.RLock()
	defer this.RUnlock()
	return this.lobby
}

func (this *LobbyHttpInterface) getUpgrader() *websocket.Upgrader {
	this.RLock()
	defer this.RUnlock()
	return this.upgrader
}
