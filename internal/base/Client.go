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
	TYPE_BOT    = "bot"
	TYPE_ENGINE = "engine"
	TYPE_PLAYER = "player"
	TYPE_VIEWER = "viewer"
)

var (
	CLIENT_TYPES = []string{TYPE_BOT, TYPE_ENGINE, TYPE_PLAYER, TYPE_VIEWER}
	CLIENT_COUNTER util.SafeCounter
)

type Client struct {
	sync.RWMutex
	Id int      `json:"id"`
	Room *Room  `json:"-"`
	Type string `json:"type"`
	Name string `json:"name"`
	Game string `json:"game"`
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
	if this.GetType() == TYPE_PLAYER {
		return // don't remove info for a player
	}
	if this.GetType() == TYPE_ENGINE && !this.GetRoom().GetIsLobby() {
		return // don't remove info for an engine in a game
	}
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
		case "action":
			handler = this.handleActionMessage
		case "invite":
			handler = this.handleInviteMessage
		case "register":
			handler = this.handleRegisterMessage
		case "room":
			handler = this.handleRoomMessage
		case "start":
			handler = this.handleStartMessage
		case "state":
			handler = this.handleStateMessage
		case "stop":
			handler = this.handleStopMessage
	}
	handler(raw, message)
}

func (this *Client) handleDefault(raw []byte, base *msg.Message) {
	message, err := msg.ParseMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	log.Warnf("No handler found for message type [%s]", message.Type)
	this.SendError(fmt.Sprintf("Invalid message type [%s]", message.Type))
}

// message handlers in alphabetical order

func (this *Client) handleActionMessage(raw []byte, base *msg.Message) {
	message, err := msg.ParseActionMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}
	message.Player = this.GetId()
	this.GetRoom().BroadcastToType(TYPE_ENGINE, message)
}

func (this *Client) handleInviteMessage(raw []byte, base *msg.Message) {
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

	engine := room.GetEngine()
	if engine != nil && client.GetType() == TYPE_BOT && engine.GetName() != client.GetName() {
		this.SendError(fmt.Sprintf("Bot [%s] is not made for [%s], but for [%s]", client.GetName(), engine.GetName(), client.GetGame()))
		return
	}

	found := room.GetClientById(client.GetId())
	if found != nil {
		this.SendError(fmt.Sprintf("Client with id [%d] is already present", message.Client))
		return
	}

	client.SendMessage(message)
}

func (this *Client) handleRegisterMessage(raw []byte, base *msg.Message) {
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

	this.GetRoom().RemoveClient(this)

	err = this.GetRoom().AddClient(this)
	if err != nil {
		log.Warnf("Could not register: [%s]", err)
		this.SendError(fmt.Sprintf("Could not register: [%s]", err))
		return
	}

	this.SetType(message.ClientType)
	this.SetName(message.Name)
	this.SetGame(message.Game)
	this.setRegistered(true)

	log.Infof("client [%s] registered on room [%s] as a [%s]", this.GetName(), this.GetRoom().GetName(), this.GetType())

	this.SendMessage(msg.NewRegisteredMessage(this.Id))
	this.GetRoom().History.SendAllToClient(this)
}

func (this *Client) handleRoomMessage(raw []byte, base *msg.Message) {
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
	this.SendMessage(msg.NewCreatedMessage(room.GetId()))
	engine.SendMessage(msg.NewInviteMessage(room.GetId(), room.GetName(), engine.GetId()))
}

func (this *Client) handleStartMessage(raw []byte, base *msg.Message) {
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

	this.GetRoom().TransformClientTypes(TYPE_BOT, TYPE_PLAYER)
	message.Players = this.GetRoom().GetClientIdsByType(TYPE_PLAYER)
	this.GetRoom().BroadcastToType(TYPE_ENGINE, message)
	this.GetRoom().SetStarted(true)
	GetLobby().TriggerUpdated()
}

func (this *Client) handleStateMessage(raw []byte, base *msg.Message) {
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

func (this *Client) handleStopMessage(raw []byte, base *msg.Message) {
	if this.GetType() != TYPE_ENGINE {
		this.SendError(fmt.Sprintf("You are not allowed to send a stop message"))
		return
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
	this.GetRoom().BroadcastToType(TYPE_PLAYER, base)
	GetLobby().TriggerUpdated()
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

func (this *Client) GetGame() string {
	this.RLock()
	defer this.RUnlock()
	return this.Game
}

func (this *Client) SetGame(game string) {
	this.Lock()
	defer this.Unlock()
	this.Game = game
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
