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
	TYPE_VIEWER = "viewer"
)

var (
	CLIENT_TYPES = []string{TYPE_BOT, TYPE_ENGINE, TYPE_VIEWER}
	CLIENT_COUNTER util.SafeCounter
)

type Client struct {
	sync.RWMutex
	Id int      `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Game string `json:"game"` // used for bots to specify which game they want to play
	lobby *Lobby
	connection *Connection
}

func NewClient(lobby *Lobby, connection *Connection) *Client {
	client :=  &Client {
		Id: CLIENT_COUNTER.GetNext(),
		lobby: lobby,
		connection: nil,
	}
	client.SetConnection(connection)
	return client
}

func (this *Client) Merge(client *Client) {
	this.getLobby().RemoveClient(client)
	*client = *this
}



// communication related stuff

func (this *Client) SendMessage(message interface{}) {
	go func() {
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
	}()
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
		case "game":
			handler = this.handleGameMessage
		case "invite":
			handler = this.handleInviteMessage
		case "join":
			handler = this.handleJoinMessage
		case "leave":
			handler = this.handleLeaveMessage
		case "register":
			handler = this.handleRegisterMessage
		case "start":
			handler = this.handleStartMessage
		case "state":
			handler = this.handleStateMessage
		case "stop":
			handler = this.handleStopMessage
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

// message handlers in alphabetical order

func (this *Client) handleActionMessage(raw []byte) {
	message, err := msg.ParseActionMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	game := this.getLobby().GetGameById(message.Game)
	if game == nil {
		this.SendError(fmt.Sprintf("game not found: [%d]", message.Game))
		return
	}

	game.HandleActionMessage(message)
}

func (this *Client) handleGameMessage(raw []byte) {
	message, err := msg.ParseGameMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	engine := this.getLobby().GetClientById(message.Engine)
	if engine == nil {
		this.SendError(fmt.Sprintf("Could not find engine with id [%d]", message.Engine))
		return
	}

	game := NewGame(message.Name, engine)
	this.getLobby().AddGame(game)
	this.SendMessage(msg.NewCreatedMessage(game.GetId()))
}

func (this *Client) handleInviteMessage(raw []byte) {
	message, err := msg.ParseInviteMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	game := this.getLobby().GetGameById(message.Game)
	if game == nil {
		this.SendError(fmt.Sprintf("Game [%d] not found", message.Game))
		return
	}

	bot := this.getLobby().GetClientById(message.Bot)
	if bot == nil {
		this.SendError(fmt.Sprintf("Bot [%d] not found", message.Bot))
		return
	}
	if bot.GetType() != TYPE_BOT {
		this.SendError(fmt.Sprintf("Client [%d] is not a bot", message.Bot))
		return
	}

	game.AddPlayer(bot)
	this.getLobby().TriggerUpdated()
}

func (this *Client) handleJoinMessage(raw []byte) {
	message, err := msg.ParseInviteMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	game := this.getLobby().GetGameById(message.Game)
	if game == nil {
		this.SendError(fmt.Sprintf("Game [%d] not found", message.Game))
		return
	}

	game.AddClient(this)
	game.GetHistory().SendAllToClient(this)
	this.getLobby().TriggerUpdated()
}

func (this *Client) handleLeaveMessage(raw []byte) {
	message, err := msg.ParseLeaveMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	game := this.getLobby().GetGameById(message.Game)
	if game == nil {
		this.SendError(fmt.Sprintf("Game [%d] not found", message.Game))
		return
	}

	game.RemoveClient(this)
	this.getLobby().TriggerUpdated()
}

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

	this.setType(message.ClientType)
	this.setName(message.Name)
	this.setGame(message.Game)

	log.Infof("client [%s] registered as a [%s]", this.GetName(), this.GetType())

	duplicate := this.getLobby().FindDuplicateUnconnectedClient(this)
	if duplicate != nil {
		log.Infof("Detected [%s] previously connected, merging", this.GetName())
		this.Merge(duplicate)
	}

	this.SendMessage(msg.NewRegisteredMessage(this.Id))
	this.getLobby().TriggerUpdated()
}

func (this *Client) handleStartMessage(raw []byte) {
	message, err := msg.ParseStartMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	game := this.getLobby().GetGameById(message.Game)
	if game == nil {
		this.SendError(fmt.Sprintf("Game [%d] not found", message.Game))
		return
	}

	err = game.Start()
	if err != nil {
		this.SendError(fmt.Sprintf("Could not start game [%d]: [%s]", message.Game, err))
		return
	}

	this.getLobby().TriggerUpdated()
}

func (this *Client) handleStateMessage(raw []byte) {
	if this.GetType() != TYPE_ENGINE {
		this.SendError(fmt.Sprintf("You are not allowed to send a state message"))
		return
	}

	message, err := msg.ParseStateMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message : [%s] : [%s]", err, raw))
		return
	}

	game := this.getLobby().GetGameById(message.Game)
	if game == nil {
		this.SendError(fmt.Sprintf("Game [%d] not found", message.Game))
		return
	}

	game.HandleStateMessage(message)
}

func (this *Client) handleStopMessage(raw []byte) {
	if this.GetType() != TYPE_ENGINE {
		this.SendError(fmt.Sprintf("You are not allowed to send a stop message"))
		return
	}

	message, err := msg.ParseStopMessage(raw)
	if err != nil {
		this.SendError(fmt.Sprintf("Could not parse message: [%s]", raw))
		return
	}

	game := this.getLobby().GetGameById(message.Game)
	if game == nil {
		this.SendError(fmt.Sprintf("Game [%d] not found", message.Game))
		return
	}

	err = game.Stop()
	if err != nil {
		this.SendError(fmt.Sprintf("Could not stop game [%d]: [%s]", message.Game, err))
		return
	}

	this.getLobby().TriggerUpdated()
}



// getters and setters

func (this *Client) GetId() int {
	this.RLock()
	defer this.RUnlock()
	return this.Id
}

func (this *Client) GetType() string {
	this.RLock()
	defer this.RUnlock()
	return this.Type
}

func (this *Client) setType(Type string) {
	this.Lock()
	defer this.Unlock()
	this.Type = Type
}

func (this *Client) GetName() string {
	this.RLock()
	defer this.RUnlock()
	return this.Name
}

func (this *Client) setName(name string) {
	this.Lock()
	defer this.Unlock()
	this.Name = name
}

func (this *Client) GetGame() string {
	this.RLock()
	defer this.RUnlock()
	return this.Game
}

func (this *Client) setGame(game string) {
	this.Lock()
	defer this.Unlock()
	this.Game = game
}

func (this *Client) getLobby() *Lobby {
	this.RLock()
	defer this.RUnlock()
	return this.lobby
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

func (this *Client) IsConnected() bool {
	this.RLock()
	defer this.RUnlock()
	return this.connection != nil
}



// lock for json marshalling

type JClient Client

func (this *Client) MarshalJSON() ([]byte, error) {
    this.RLock()
    defer this.RUnlock()
    return json.Marshal(JClient(*this))
}
