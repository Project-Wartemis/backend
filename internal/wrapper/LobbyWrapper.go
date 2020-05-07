package wrapper

import (
	"net/http"
	"github.com/Project-Wartemis/pw-backend/internal/base"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

type LobbyWrapper struct {
}

func NewLobbyWrapper() *LobbyWrapper {
	return &LobbyWrapper {}
}

func (this *LobbyWrapper) GetBots(writer http.ResponseWriter, request *http.Request) {
	bots := []*base.Client{}
	for _, client := range base.GetLobby().Clients {
		if client.IsBot {
			bots = append(bots, client)
		}
	}
	util.WriteJson(writer, bots)
}

func (this *LobbyWrapper) NewConnection(writer http.ResponseWriter, request *http.Request) {
	client := base.GetLobby().CreateAndAddClient()
	defer base.GetLobby().RemoveClient(client)

	util.SetupWebSocket(writer, request, client.SetConnection, client.HandleMessage)
}
