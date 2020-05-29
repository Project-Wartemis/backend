package base

import (
	"encoding/json"
	"github.com/google/uuid"
	sync "github.com/sasha-s/go-deadlock"
)

type Player struct {
	sync.RWMutex
	Id int `json:"id"`
	Client *Client `json:"client"`
	key string
}

func NewPlayer(id int, client *Client) *Player {
	return &Player {
		Id: id,
		Client: client,
		key: uuid.New().String(),
	}
}



// getters and setters

func (this *Player) GetId() int {
	this.RLock()
	defer this.RUnlock()
	return this.Id
}

func (this *Player) GetClient() *Client {
	this.RLock()
	defer this.RUnlock()
	return this.Client
}

func (this *Player) GetKey() string {
	this.RLock()
	defer this.RUnlock()
	return this.key
}



// lock for json marshalling

type JPlayer Player

func (this *Player) MarshalJSON() ([]byte, error) {
    this.RLock()
    defer this.RUnlock()
    return json.Marshal(JPlayer(*this))
}
