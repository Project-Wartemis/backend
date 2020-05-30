package base

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/message"
)

type Lobby struct {
	Room
	Games []*Game `json:"games"`
	gamesById map[int]*Game
}

func NewLobby() *Lobby {
	room := NewRoom("lobby")
	return &Lobby {
		Room: *room,
		Games: []*Game{},
		gamesById: map[int]*Game{},
	}
}

func (this *Lobby) HandleConnect(connection *Connection) {
	client := NewClient(this, connection)
	this.AddClient(client)
	log.Info("Added a new client")
	client.SendMessage(message.NewConnectedMessage())
}

func (this *Lobby) HandleDisconnect(client *Client) {
	if client.GetType() == TYPE_VIEWER {
		this.RemoveClient(client)
		this.RLock()
		for _,game := range this.Games {
			game.RemoveClient(client)
		}
		this.RUnlock()
	}
}

func (this *Lobby) HandleReconnect(new *Client, old *Client) {
	log.Infof("Reconnecting [%s]", new.GetName())
	old.Transfer(new)
	this.RemoveClient(new)
	this.RLock()
	defer this.RUnlock()
	for _,game := range this.Games {
		game.HandleReconnect(old)
	}
}



// communication related stuff

func (this *Lobby) SendMessage(clientId int, message interface{}) {
	client := this.GetClientById(clientId)
	if client == nil {
		log.Warnf("Tried sending a message to client [%d], but not found in [%s]", clientId, this.GetName())
		return
	}
	client.SendMessage(message)
}

func (this *Lobby) TriggerUpdated() {
	this.BroadcastToType(TYPE_VIEWER, message.NewLobbyMessage(this))
}



// getters and setters

func (this *Lobby) AddGame(game *Game) {
	log.Infof("Adding game [%s]", game.GetName())

	this.setGameById(game.GetId(), game)

	this.Lock()
	defer this.Unlock()
	this.Games = append(this.Games, game)

	go this.TriggerUpdated()
}

func (this *Lobby) RemoveGame(game *Game) {
	log.Infof("Removing game [%s]", game.GetName())

	this.removeGameById(game.GetId())

	this.Lock()
	defer this.Unlock()
	for i,r := range this.Games {
		if r.Id == game.Id {
			this.Games[i] = this.Games[len(this.Games)-1] // copy last element to index i
			this.Games[len(this.Games)-1] = nil           // erase last element
			this.Games = this.Games[:len(this.Games)-1]   // truncate slice
		}
	}

	go this.TriggerUpdated()
}

func (this *Lobby) GetGameById(id int) *Game {
	this.RLock()
	defer this.RUnlock()
	return this.gamesById[id]
}

func (this *Lobby) setGameById(id int, game *Game) {
	this.Lock()
	defer this.Unlock()
	this.gamesById[id] = game
}

func (this *Lobby) removeGameById(id int) {
	this.Lock()
	defer this.Unlock()
	delete(this.gamesById, id)
}



// lock for json marshalling

type JLobby Lobby

func (this *Lobby) MarshalJSON() ([]byte, error) {
    this.RLock()
    defer this.RUnlock()
    return json.Marshal(JLobby(*this))
}
