package base

import (
	log "github.com/sirupsen/logrus"
)

type Room struct {
	Name string
	Clients []*Client
	clientsByKey map[string]*Client
	clientsByName map[string]*Client
}

func (this *Room) AddClient(client *Client) {
	this.Clients = append(this.Clients, client)
	log.Info("Added client to room [%s]", this.Name)
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
