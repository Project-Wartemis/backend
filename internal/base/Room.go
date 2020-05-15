package base

import (
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	Name string                `json:"name"`
	Key string                 `json:"key"`
	Clients map[string][]*Client `json:"clients"`
	clientsByKey map[string]*Client
}

func NewRoom(name string) *Room {
	return &Room {
		Name: name,
		Key: uuid.New().String(),
		Clients: map[string][]*Client{},
		clientsByKey: map[string]*Client{},
	}
}

func (this *Room) GetClientByKey(key string) *Client {
	return this.clientsByKey[key]
}

func (this *Room) CreateAndAddClient() *Client {
	client := NewClient(this)
	log.Infof("Added a new client to room [%s]", this.Name)
	return client
}

func (this *Room) AddClient(client *Client) error {
	if client.Key != "" {
		if _, found := this.clientsByKey[client.Key]; found {
			return errors.New("key already registered")
		}
		this.clientsByKey[client.Key] = client
	}
	if client.Type != "" {
		this.Clients[client.Type] = append(this.Clients[client.Type], client)
	}
	return nil
}

func (this *Room) RemoveClient(client *Client) {
	if client.Key != "" {
		delete(this.clientsByKey, client.Key)
		log.Infof("Removing client [%s] from room [%s]", client.Name, this.Name)
	}
	if client.Type != "" {
		this.Clients[client.Type] = this.removeClientFromList(client, this.Clients[client.Type])
	}
}

func (this *Room) removeClientFromList(client *Client, list []*Client) []*Client {
	for i,c := range list {
		if c != client {
			continue
		}
		list[i] = list[len(list)-1] // copy last element to index i
		list[len(list)-1] = nil     // erase last element
		list = list[:len(list)-1]   // truncate slice
	}
	return list
}

func (this *Room) SendMessage(key string, message interface{}) {
	if client, found := this.clientsByKey[key]; found {
		client.SendMessage(message)
	}
}

func (this *Room) Broadcast(client *Client, message interface{}) {
	for _,c := range this.clientsByKey {
		if c != client {
			c.SendMessage(message)
		}
	}
}
