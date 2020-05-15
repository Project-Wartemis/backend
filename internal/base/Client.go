package base

import (
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/message"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

const (
	TYPE_VIEWER = "viewer"
	TYPE_BOT    = "bot"
	TYPE_ENGINE = "engine"
)

var CLIENT_TYPES = []string{TYPE_VIEWER, TYPE_BOT, TYPE_ENGINE}

type Client struct {
	Room *Room  `json:"-"`
	Type string `json:type`
	Name string `json:"name"`
	Key string  `json:"key"`
	isRegistered bool
	connection *websocket.Conn
}

func NewClient(room *Room) *Client {
	return &Client {
		Room: room,
		Type: "",
		isRegistered: false,
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
	this.SendMessage(struct {
		Type string    `json:"type"`
		Message string `json:"message"`
	}{
		Type: "error",
		Message: message,
	})
}

func (this *Client) HandleMessage(raw []byte) {
	log.Debugf("got message: %s", raw)
	message, err := message.ParseMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}
	handler := this.handleDefault
	switch message.Type {
		case "echo":
			handler = this.handleEchoMessage
		case "register":
			handler = this.handleRegisterMessage
		case "gamestate":
			handler = this.handleGamestateMessage
		case "request_move":
			handler = this.handleMoveRequestMessage
	}
	handler(raw)
}

func (this *Client) handleDefault(raw []byte) {
	message, err := message.ParseMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}

	log.Warnf("No handler found for message type [%s]", message.Type)
	this.SendError(fmt.Sprintf("Invalid message type [%s]", message.Type))
}

func (this *Client) handleEchoMessage(raw []byte) {
	message, err := message.ParseEchoMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}
	log.Debugf("Echoing value [%s]", message.Value)
	this.SendMessage(struct {
		Value string `json:"value"`
	}{
		Value: message.Value,
	})
}

func (this *Client) handleRegisterMessage(raw []byte) {
	message, err := message.ParseRegisterMessage(raw)
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
	this.Key = message.Key
	this.isRegistered = false

	err = this.Room.AddClient(this);
	if err != nil {
		log.Warn("Could not register")
		this.SendError(fmt.Sprintf("Could not register: [%s]", err))
		return
	}

	this.isRegistered = true
	log.Infof("client [%s] registered on room [%s] as a [%s]", this.Name, this.Room.Name, this.Type)
}

func (this *Client) handleGamestateMessage(raw []byte) {
	message, err := message.ParseGamestateMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}
	GetLobby().Broadcast(this, message)
}

func (this *Client) handleMoveRequestMessage(raw []byte) {
	message, err := message.ParseMoveRequestMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw));
		return;
	}
	GetLobby().SendMessage(message.Key, message)
}
