package base

import (
	"errors"
	"sync"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	Clients []*Client
	clientsByKey map[string]*Client
	clientsByName map[string]*Client
}

// singleton lobby
var (
	lobby *Room
)

func GetLobby() *Room {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	if lobby == nil {
		lobby = &Room {
			Clients: []*Client{},
			clientsByKey: make(map[string]*Client),
			clientsByName: make(map[string]*Client),
		}
	}
	return lobby
}

func (this *Room) CreateAndAddClient() *Client {
	client := NewClient()
	this.Clients = append(this.Clients, client)
	log.Info("Added client")
	return client
}

func (this *Room) Register(client *Client, name string, key string) error {
	if this != GetLobby() {
		log.Errorf("Can only register on the lobby")
		return errors.New("internal server error")
	}
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

func (this *Room) RemoveClient(client *Client) {
	for i,c := range this.Clients {
		if c != client {
			continue
		}
		this.Clients[i] = this.Clients[len(this.Clients)-1] // copy last element to index i
		this.Clients[len(this.Clients)-1] = nil             // erase last element
		this.Clients = this.Clients[:len(this.Clients)-1]   // truncate slice
		log.Infof("Removed client [%s]", client.Name)
		return
	}
}

func (this *Room) SendMessage(key string, message interface{}) {
	for _,c := range this.Clients {
		if(c.Key == key) {
			c.SendMessage(message)
		}
	}
}

func (this *Room) Broadcast(client *Client, message interface{}) {
	for _,c := range this.Clients {
		if(c != client) {
			c.SendMessage(message)
		}
	}
}
