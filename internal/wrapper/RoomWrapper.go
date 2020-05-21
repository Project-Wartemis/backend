package wrapper

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/Project-Wartemis/pw-backend/internal/base"
	"github.com/Project-Wartemis/pw-backend/internal/message"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

type RoomWrapper struct {}

func NewRoomWrapper() *RoomWrapper {
	return &RoomWrapper {}
}

func (this *RoomWrapper) NewConnection(writer http.ResponseWriter, request *http.Request) {
	roomIdText := mux.Vars(request)["room"]

	roomId, err := strconv.Atoi(roomIdText)
	if err != nil {
		util.WriteStatus(writer, http.StatusNotFound, fmt.Sprintf("Could not parse [%s] to id", roomIdText))
		return
	}

	room := base.GetLobby().GetRoomById(roomId)
	if room == nil {
		util.WriteStatus(writer, http.StatusNotFound, fmt.Sprintf("Could not find room for id [%d]", roomId))
		return
	}

	this.newConnection(room, writer, request)
}

func (this *RoomWrapper) newConnection(room *base.Room, writer http.ResponseWriter, request *http.Request) {
	client := room.CreateAndAddClient()
	defer client.HandleDisconnect()

	postInit := func(conn *websocket.Conn) {
		connection := base.NewConnection(conn)
		client.SetConnection(connection)
		client.SendMessage(message.NewConnectedMessage())
	}

	util.SetupWebSocket(writer, request, postInit, client.HandleMessage)
}
