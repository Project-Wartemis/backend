package base

import (
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/message"
)

type Client struct {
	Name string `json:"name"`
	Key string  `json:"key"`
	IsBot bool  `json:"-"`
	connection *websocket.Conn
}

func NewClient() *Client {
	return &Client {}
}

func (this *Client) SetConnection(connection *websocket.Conn) {
	this.connection = connection
}

func (this *Client) SendMessage(message interface{}) {
	log.Debugf("Sending message to [%s]: [%s]", this.Name, message)
	text, err := json.Marshal(message)
	if err != nil {
		log.Error("Unexpected error while parsing message to json: [%s]", message)
		return
	}
	err = this.connection.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		log.Error("Unexpected error while sending message to [%s] : [%s]", this.Name, message)
		return
	}
}

func (this *Client) SendError(message string) {
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
	message := message.ParseMessage(raw)
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
	message := message.ParseMessage(raw)
	log.Warnf("No handler found for message type [%s]", message.Type)
	this.SendError(fmt.Sprintf("Invalid message type [%s]", message.Type))
}

func (this *Client) handleEchoMessage(raw []byte) {
	message := message.ParseEchoMessage(raw)
	log.Debugf("Echoing value [%s]", message.Value)
	this.SendMessage(struct {
		Value string `json:"value"`
	}{
		Value: message.Value,
	})
}

func (this *Client) handleRegisterMessage(raw []byte) {
	message := message.ParseRegisterMessage(raw)
	err := GetLobby().Register(this, message.Name, message.Key)
	if err != nil {
		log.Warn("Could not register")
		this.SendError(fmt.Sprintf("Could not register: [%s]", err))
	}
}

func (this *Client) handleGamestateMessage(raw []byte) {
	message := message.ParseGamestateMessage(raw)
	GetLobby().Broadcast(this, message)
}

func (this *Client) handleMoveRequestMessage(raw []byte) {
	message := message.ParseMoveRequestMessage(raw)
	GetLobby().SendMessage(message.Key, message)
}
