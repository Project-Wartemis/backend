package base

import (
	"encoding/json"
	"fmt"
	sync "github.com/sasha-s/go-deadlock"
	log "github.com/sirupsen/logrus"
	msg "github.com/Project-Wartemis/pw-backend/internal/message"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

const (
	TYPE_VIEWER = "viewer"
	TYPE_BOT    = "bot"
	TYPE_ENGINE = "engine"
)

var (
	CLIENT_TYPES = []string{TYPE_VIEWER, TYPE_BOT, TYPE_ENGINE}
	CLIENT_COUNTER util.SafeCounter
)

type Client struct {
	sync.RWMutex
	Id int      `json:"id"`
	Room *Room  `json:"-"`
	Type string `json:"type"`
	Name string `json:"name"`
	registered bool
	connection *Connection
}

func NewClient(room *Room) *Client {
	return &Client {
		Id: CLIENT_COUNTER.GetNext(),
		Room: room,
		Type: "",
		registered: false,
		connection: nil,
	}
}



// basic communication related stuff

func (this *Client) HandleDisconnect() {
	this.SetConnection(nil)
	this.Room.RemoveClient(this)
	GetLobby().TriggerUpdated()
}

func (this *Client) SendMessage(message interface{}) {
	log.Debugf("Sending message to [%s]: [%s]", this.GetName(), message)

	connection := this.getConnection()
	if connection == nil {
		log.Warnf("Cannot send a message to [%s] because not connected", this.GetName())
		return
	}

	err := connection.SendMessage(message)
	if err != nil {
		log.Errorf("Unexpected error while sending message to [%s] : [%s]", this.GetName(), err)
		return
	}
}

func (this *Client) SendError(message string) {
	log.Infof("Sending error message to [%s]: [%s]", this.GetName(), message)
	this.SendMessage(msg.NewErrorMessage(message))
}



// message handling

func (this *Client) HandleMessage(raw []byte) {
	log.Debugf("got message: %s", raw)
	message, err := msg.ParseMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
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
		case "stop":
			handler = this.handleStopMessage
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
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	log.Warnf("No handler found for message type [%s]", message.Type)
	this.SendError(fmt.Sprintf("Invalid message type [%s]", message.Type))
}

// message handlers in alphabetical order TODO

func (this *Client) handleRegisterMessage(raw []byte) {
	message, err := msg.ParseRegisterMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	if !util.Includes(CLIENT_TYPES, message.ClientType) {
		this.SendError(fmt.Sprintf("Could not register: [Invalid value for clientType [%s]]", message.ClientType))
		return
	}

	defer GetLobby().TriggerUpdated()

	this.GetRoom().RemoveClient(this) // TODO

	err = this.GetRoom().AddClient(this) // TODO
	if err != nil {
		log.Warnf("Could not register: [%s]", err)
		this.SendError(fmt.Sprintf("Could not register: [%s]", err))
		return
	}

	this.SetType(message.ClientType)
	this.SetName(message.Name)
	this.setRegistered(true)

	log.Infof("client [%s] registered on room [%s] as a [%s]", this.GetName(), this.GetRoom().GetName(), this.GetType())

	this.SendMessage(msg.NewRegisteredMessage(this.Id))
	this.GetRoom().History.SendAllToClient(this)
}

func (this *Client) handleRoomMessage(raw []byte) {
	message, err := msg.ParseRoomMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	engine := GetLobby().GetClientById(message.Engine)
	if engine == nil {
		this.SendError(fmt.Sprintf("Could not find engine with id [%d]", message.Engine))
		return
	}

	room := NewRoom(message.Name, false)
	GetLobby().AddRoom(room)
	this.SendMessage(msg.NewInviteMessage(room.GetId(), room.GetName(), this.GetId()))
	engine.SendMessage(msg.NewInviteMessage(room.GetId(), room.GetName(), engine.GetId()))
}

func (this *Client) handleInviteMessage(raw []byte) {
	message, err := msg.ParseInviteMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	room := GetLobby().GetRoomById(message.Room)
	if room == nil {
		this.SendError(fmt.Sprintf("Could not find room with id [%d]", message.Room))
		return
	}
	message.Name = room.Name

	client := GetLobby().GetClientById(message.Client)
	if client == nil {
		this.SendError(fmt.Sprintf("Could not find player with id [%d]", message.Client))
		return
	}

	found := room.GetClientById(client.GetId())
	if found != nil {
		this.SendError(fmt.Sprintf("Client with id [%d] is already present", message.Client))
		return
	}

	client.SendMessage(message)
}

func (this *Client) handleStartMessage(raw []byte) {
	if this.GetRoom().GetIsLobby() {
		this.SendError(fmt.Sprintf("You cannot start a game in the lobby..."))
		return
	}
	if this.GetRoom().GetStarted() {
		this.SendError(fmt.Sprintf("This game has already started"))
		return
	}

	message, err := msg.ParseStartMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	message.Players = this.GetRoom().GetClientIdsByType(TYPE_BOT)
	this.GetRoom().BroadcastToType(TYPE_ENGINE, message)
	this.GetRoom().SetStarted(true)
	GetLobby().TriggerUpdated()
}

func (this *Client) handleStopMessage(raw []byte) {
	if this.GetType() != TYPE_ENGINE {
		// TODO uncomment once a proper engine is implemented
		//this.SendError(fmt.Sprintf("You are not allowed to send a stop message"))
		//return
	}
	if !this.GetRoom().GetStarted() {
		this.SendError(fmt.Sprintf("This game has not started yet"))
		return
	}
	if this.GetRoom().GetStopped() {
		this.SendError(fmt.Sprintf("This game has already stopped"))
		return
	}

	this.GetRoom().SetStopped(true)
	GetLobby().TriggerUpdated()
}

func (this *Client) handleStateMessage(raw []byte) {
	if this.GetType() != TYPE_ENGINE {
		this.SendError(fmt.Sprintf("You are not allowed to send a state message"))
		return
	}

	message, err := msg.ParseStateMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	this.GetRoom().History.Add(message)
	this.GetRoom().Broadcast(this.GetId(), message)
}

func (this *Client) handleActionMessage(raw []byte) {
	message, err := msg.ParseActionMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}
	message.Player = this.GetId()
	this.GetRoom().BroadcastToType(TYPE_ENGINE, message)
}



// getters and setters

func (this *Client) GetId() int {
	this.RLock()
	defer this.RUnlock()
	return this.Id
}

func (this *Client) SetId(id int) {
	this.Lock()
	defer this.Unlock()
	this.Id = id
}

func (this *Client) GetRoom() *Room {
	this.RLock()
	defer this.RUnlock()
	return this.Room
}

// SetRoom not implemented, it should not update

func (this *Client) GetType() string {
	this.RLock()
	defer this.RUnlock()
	return this.Type
}

func (this *Client) SetType(Type string) {
	this.Lock()
	defer this.Unlock()
	this.Type = Type
}

func (this *Client) GetName() string {
	this.RLock()
	defer this.RUnlock()
	return this.Name
}

func (this *Client) SetName(name string) {
	this.Lock()
	defer this.Unlock()
	this.Name = name
}

func (this *Client) getRegistered() bool {
	this.RLock()
	defer this.RUnlock()
	return this.registered
}

func (this *Client) setRegistered(registered bool) {
	this.Lock()
	defer this.Unlock()
	this.registered = registered
}

func (this *Client) getConnection() *Connection {
	this.RLock()
	defer this.RUnlock()
	return this.connection
}

func (this *Client) SetConnection(connection *Connection) {
	if this.getConnection() != nil {
		this.getConnection().StopPinging()
	}

	this.Lock()
	defer this.Unlock()

	this.connection = connection
	if connection != nil {
		go connection.StartPinging()
	}
}



// lock for json marshalling

type JClient Client

func (this *Client) MarshalJSON() ([]byte, error) {
    this.RLock()
    defer this.RUnlock()
    return json.Marshal(JClient(*this))
}
