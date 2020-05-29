package base

import (
	"encoding/json"
	sync "github.com/sasha-s/go-deadlock"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	sync.RWMutex
	Name string         `json:"name"`
	Clients []*Client `json:"clients"`
	clientsById map[int]*Client
}

func NewRoom(name string) *Room {
	return &Room {
		Name: name,
		Clients: []*Client{},
		clientsById: map[int]*Client{},
	}
}



// basic communication related stuff

func (this *Room) Broadcast(senderId int, message interface{}) {
	this.RLock()
	defer this.RUnlock()
	for _,client := range this.Clients {
		if client.GetId() != senderId {
			client.SendMessage(message)
		}
	}
}

func (this *Room) BroadcastToType(Type string, message interface{}) {
	this.RLock()
	defer this.RUnlock()
	for _,client := range this.Clients {
		if client.GetType() == Type {
			client.SendMessage(message)
		}
	}
}



// getters and setters

func (this *Room) GetName() string {
	this.RLock()
	defer this.RUnlock()
	return this.Name
}

func (this *Room) AddClient(client *Client) {
	this.setClientById(client.GetId(), client)

	this.Lock()
	defer this.Unlock()
	this.Clients = append(this.Clients, client)
}

func (this *Room) RemoveClient(client *Client) {
	log.Infof("Removing client [%s] from room [%s]", client.GetName(), this.GetName())

	this.Lock()
	defer this.Unlock()
	for i,c := range this.Clients {
		if c.GetId() == client.GetId() {
			this.Clients[i] = this.Clients[len(this.Clients)-1] // copy last element to index i
			this.Clients[len(this.Clients)-1] = nil             // erase last element
			this.Clients = this.Clients[:len(this.Clients)-1]   // truncate slice
			return
		}
	}
}

func (this *Room) FindDuplicateUnconnectedClient(client *Client) *Client {
	this.RLock()
	defer this.RUnlock()
	for _,c := range this.Clients {
		if c.GetType() != client.GetType() {
			continue
		}
		if c.GetName() != client.GetName() {
			continue
		}
		if c.IsConnected() {
			continue
		}
		return c
	}
	return nil
}

func (this *Room) GetClientById(id int) *Client {
	this.RLock()
	defer this.RUnlock()
	return this.clientsById[id]
}

func (this *Room) setClientById(id int, client *Client) {
	this.Lock()
	defer this.Unlock()
	this.clientsById[id] = client
}

func (this *Room) removeClientById(id int) {
	this.Lock()
	defer this.Unlock()
	delete(this.clientsById, id)
}



// lock for json marshalling

type JRoom Room

func (this *Room) MarshalJSON() ([]byte, error) {
    this.RLock()
    defer this.RUnlock()
    return json.Marshal(JRoom(*this))
}
