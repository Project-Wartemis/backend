package base

import (
	"fmt"
	"sync"
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	msg "github.com/Project-Wartemis/pw-backend/internal/message"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

const (
	TYPE_VIEWER = "viewer"
	TYPE_BOT    = "bot"
	TYPE_ENGINE = "engine"
)

var CLIENT_TYPES = []string{TYPE_VIEWER, TYPE_BOT, TYPE_ENGINE}

var CLIENT_COUNTER = 0

func getNextClientId() int {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	CLIENT_COUNTER++
	return CLIENT_COUNTER
}

type Client struct {
	Id int      `json:"id"`
	Room *Room  `json:"-"`
	Type string `json:"-"`
	Name string `json:"name"`
	isRegistered bool
	connection *websocket.Conn
}

func NewClient(room *Room) *Client {
	return &Client {
		Id: getNextClientId(),
		Room: room,
		Type: "",
		isRegistered: false,
		connection: nil,
	}
}

func (this *Client) SetConnection(connection *websocket.Conn) {
	this.connection = connection
}

func (this *Client) SendMessage(message interface{}) {
	log.Debugf("Sending message to [%s]: [%s]", this.Name, message)
	text, err := json.Marshal(message)
	if err != nil {
		log.Errorf("Unexpected error while parsing message to json: [%s]", message)
		return
	}
	err = this.connection.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		log.Errorf("Unexpected error while sending message to [%s] : [%s]", this.Name, message)
		return
	}
}

func (this *Client) SendError(message string) {
	log.Infof("Sending error message to [%s]: [%s]", this.Name, message)
	this.SendMessage(msg.NewErrorMessage(message))
}

func (this *Client) HandleMessage(raw []byte) {
	log.Debugf("got message: %s", raw)
	message, err := msg.ParseMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}
	handler := this.handleDefault
	switch message.Type {
		case "register":
			handler = this.handleRegisterMessage
		case "room":
			handler = this.handleRoomMessage
		case "invite":
			handler = this.handleInviteMessage
		case "start":
			handler = this.handleStartMessage
		case "state":
			handler = this.handleStateMessage
		case "action":
			handler = this.handleActionMessage
	}
	handler(raw)
}

func (this *Client) handleDefault(raw []byte) {
	message, err := msg.ParseMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}

	log.Warnf("No handler found for message type [%s]", message.Type)
	this.SendError(fmt.Sprintf("Invalid message type [%s]", message.Type))
}

func (this *Client) handleRegisterMessage(raw []byte) {
	message, err := msg.ParseRegisterMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}

	if !util.Includes(CLIENT_TYPES, message.ClientType) {
		this.SendError(fmt.Sprintf("Could not register: [Invalid value for clientType [%s]]", message.ClientType))
		return
	}

	this.Room.RemoveClient(this);

	this.Type = message.ClientType
	this.Name = message.Name
	this.isRegistered = false

	err = this.Room.AddClient(this);
	if err != nil {
		log.Warn("Could not register")
		this.SendError(fmt.Sprintf("Could not register: [%s]", err))
		return
	}

	this.isRegistered = true
	log.Infof("client [%s] registered on room [%s] as a [%s]", this.Name, this.Room.Name, this.Type)

	this.SendMessage(msg.NewRegisteredMessage(this.Id))
}

func (this *Client) handleRoomMessage(raw []byte) {
	message, err := msg.ParseRoomMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}

	engine := GetLobby().GetClientById(message.Engine)
	if engine == nil {
		this.SendError(fmt.Sprintf("Could not find engine with id [%d]", message.Engine))
		return
	}

	room := GetLobby().CreateAndAddRoom(message.Name)

	this.SendMessage(msg.NewInviteMessage(room.Id, this.Id))
	engine.SendMessage(msg.NewInviteMessage(room.Id, engine.Id))
}

func (this *Client) handleInviteMessage(raw []byte) {
	message, err := msg.ParseInviteMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}

	room := GetLobby().GetRoomById(message.Room)
	if room == nil {
		this.SendError(fmt.Sprintf("Could not find room with id [%d]", message.Room))
		return
	}

	client := GetLobby().GetClientById(message.Client)
	if client == nil {
		this.SendError(fmt.Sprintf("Could not find player with id [%d]", message.Client))
		return
	}

	found := room.GetClientById(client.Id)
	if found != nil {
		this.SendError(fmt.Sprintf("Client with id [%d] is already present", message.Client))
		return
	}

	client.SendMessage(message)
}

func (this *Client) handleStartMessage(raw []byte) {
	message, err := msg.ParseStartMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}

	message.Players = this.Room.GetBotIds()
	this.Room.SendMessageToEngine(message)
}

func (this *Client) handleStateMessage(raw []byte) {
	message, err := msg.ParseStateMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}
	this.Room.Broadcast(this, message)
}

func (this *Client) handleActionMessage(raw []byte) {
	message, err := msg.ParseActionMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}
	this.Room.SendMessageToEngine(message)
}
