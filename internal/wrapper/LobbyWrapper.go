package wrapper

import (
	"net/http"
	"github.com/Project-Wartemis/pw-backend/internal/base"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

type LobbyWrapper struct {
	roomWrapper *RoomWrapper
}

func NewLobbyWrapper(roomWrapper *RoomWrapper) *LobbyWrapper {
	return &LobbyWrapper {
		roomWrapper: roomWrapper,
	}
}

func (this *LobbyWrapper) GetLobby(writer http.ResponseWriter, request *http.Request) {
	util.WriteJson(writer, base.GetLobby())
}

func (this *LobbyWrapper) NewConnection(writer http.ResponseWriter, request *http.Request) {
	this.roomWrapper.newConnection(&(base.GetLobby().Room), writer, request)
}
