package base

import (
	"errors"
	"sync"
	log "github.com/sirupsen/logrus"
)

type Lobby struct {
	Room
	Rooms []*Room
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

	room := &Room {
		Clients: []*Client{},
		clientsByKey: make(map[string]*Client),
		clientsByName: make(map[string]*Client),
	}
	lobby = &Lobby {
		Room: *room,
		Rooms: []*Room{},
	}

	return lobby
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
	if _, found := this.clientsByName[name]; found {
		return errors.New("name already registered")
	}

	client.Name = name
	client.IsBot = true
	client.Key = key
	this.clientsByKey[key] = client
	this.clientsByName[name] = client

	log.Infof("client [%s] registered with key [%s]", name, key)
	return nil
}
