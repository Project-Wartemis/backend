package base

import (
	"errors"
	"sync"
	log "github.com/sirupsen/logrus"
)

type Lobby struct {
	Room
	Rooms []*Room `json:"rooms"`
	roomsByKey map[string]*Room
}

var (
	lobby *Lobby
)

func GetLobby() *Lobby {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	if lobby != nil {
		return lobby
	}

	room := NewRoom()
	lobby = &Lobby {
		Room: *room,
		Rooms: []*Room{},
		roomsByKey: map[string]*Room{},
	}

	return lobby
}

func (this *Lobby) GetRoomByKey(key string) *Room {
	return this.roomsByKey[key]
}

func (this *Lobby) AddRoom(room *Room) *Room {
	result := NewRoom()
	result.Name = room.Name
	this.Rooms = append(this.Rooms, result)
	this.roomsByKey[result.Key] = result
	return result
}

func (this *Lobby) CreateAndAddClient() *Client {
	client := NewClient()
	this.AddClient(client)
	return client
}

func (this *Lobby) Register(client *Client, name string, key string) error {
	if _, found := this.clientsByKey[key]; found {
		return errors.New("key already registered")
	}

	client.Name = name
	client.IsBot = true
	client.Key = key
	this.clientsByKey[key] = client

	log.Infof("client [%s] registered with key [%s]", name, key)
	return nil
}
