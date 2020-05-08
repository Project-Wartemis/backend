package wrapper

import (
	"encoding/json"
	"net/http"
	"github.com/Project-Wartemis/pw-backend/internal/base"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

type LobbyWrapper struct {}

func NewLobbyWrapper() *LobbyWrapper {
	return &LobbyWrapper {}
}

func (this *LobbyWrapper) GetLobby(writer http.ResponseWriter, request *http.Request) {
	util.WriteJson(writer, base.GetLobby())
}

func (this *LobbyWrapper) NewConnection(writer http.ResponseWriter, request *http.Request) {
	client := base.GetLobby().CreateAndAddClient()
	defer base.GetLobby().RemoveClient(client)

	util.SetupWebSocket(writer, request, client.SetConnection, client.HandleMessage)
}

func (this *LobbyWrapper) NewRoom(writer http.ResponseWriter, request *http.Request) {
	room := new(base.Room)
	err := json.NewDecoder(request.Body).Decode(room)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	room = base.GetLobby().AddRoom(room)
	util.WriteJson(writer, room)
}
