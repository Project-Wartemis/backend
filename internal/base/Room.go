package base

import (
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	Name string          `json:"name"`
	Key string           `json:"key"`
	Bots []*Client       `json:"bots"`
	Spectators []*Client `json:"spectators"`
	clientsByKey map[string]*Client
}

func NewRoom(name string) *Room {
	return &Room {
		Name: name,
		Key: uuid.New().String(),
		Bots: []*Client{},
		Spectators: []*Client{},
		clientsByKey: map[string]*Client{},
	}
}

func (this *Room) GetClientByKey(key string) *Client {
	return this.clientsByKey[key]
}

func (this *Room) CreateAndAddClient() *Client {
	client := NewClient(this)
	this.Spectators = append(this.Spectators, client)
	log.Infof("Added client [%s] to room [%s]", client.Name, this.Name)
	return client
}

func (this *Room) RemoveClient(client *Client) {
	this.removeClientFromList(client, &this.Bots)
	this.removeClientFromList(client, &this.Spectators)
	delete(this.clientsByKey, client.Key)
	log.Infof("Removed client [%s] from room [%s]", client.Name, this.Name)
}

func (this *Room) removeClientFromList(client *Client, list *[]*Client) {
	for i,c := range *list {
		if c != client {
			continue
		}
		(*list)[i] = (*list)[len(*list)-1] // copy last element to index i
		(*list)[len(*list)-1] = nil        // erase last element
		*list = (*list)[:len(*list)-1]     // truncate slice
	}
}

func (this *Room) SendMessage(key string, message interface{}) {
	if client, found := this.clientsByKey[key]; found {
		client.SendMessage(message)
	}
}

func (this *Room) Broadcast(client *Client, message interface{}) {
	for _,c := range this.clientsByKey {
		if(c != client) {
			c.SendMessage(message)
		}
	}
}

func (this *Room) Register(client *Client, name string, key string) error {
	if _, found := this.clientsByKey[key]; found {
		return errors.New("key already registered")
	}


	client.Name = name
	client.IsBot = true
	client.Key = key

	this.removeClientFromList(client, &this.Spectators)
	this.Bots = append(this.Bots, client)
	this.clientsByKey[key] = client

	log.Infof("client [%s] registered with key [%s]", name, key)
	return nil
}
