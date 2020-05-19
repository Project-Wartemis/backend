package base

import (
	"errors"
	"sync"
	log "github.com/sirupsen/logrus"
)

var ROOM_COUNTER = 0

func getNextRoomId() int {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	ROOM_COUNTER++
	return ROOM_COUNTER
}

type Room struct {
	Id int                       `json:"id"`
	Name string                  `json:"name"`
	Clients map[string][]*Client `json:"clients"`
	engine *Client
	clientsById map[int]*Client
}

func NewRoom(name string) *Room {
	return &Room {
		Id: getNextRoomId(),
		Name: name,
		Clients: map[string][]*Client{},
		engine: nil,
		clientsById: map[int]*Client{},
	}
}

func (this *Room) GetClientById(id int) *Client {
	return this.clientsById[id]
}

func (this *Room) GetBotIds() []int {
	result := []int{}
	for id := range this.clientsById {
		result = append(result, id)
	}
	return result
}

func (this *Room) CreateAndAddClient() *Client {
	client := NewClient(this)
	this.AddClient(client)
	log.Infof("Added a new client to room [%s]", this.Name)
	return client
}

func (this *Room) AddClient(client *Client) error {
	if client.Type == TYPE_ENGINE && this.engine != nil {
		return errors.New("engine already registered on this room")
	}
	this.clientsById[client.Id] = client
	if client.Type != "" {
		this.Clients[client.Type] = append(this.Clients[client.Type], client)
	}
	if client.Type == "engine" {
		this.engine = client
	}
	GetLobby().TriggerUpdated()
	return nil
}

func (this *Room) RemoveClient(client *Client) {
	log.Infof("Removing client [%s] from room [%s]", client.Name, this.Name)
	delete(this.clientsById, client.Id)
	if client.Type != "" {
		this.Clients[client.Type] = this.removeClientFromList(client, this.Clients[client.Type])
	}
	if client.Type == "engine" {
		this.engine = nil
	}
	GetLobby().TriggerUpdated()
}

func (this *Room) removeClientFromList(client *Client, list []*Client) []*Client {
	for i,c := range list {
		if c.Id != client.Id {
			continue
		}
		list[i] = list[len(list)-1] // copy last element to index i
		list[len(list)-1] = nil     // erase last element
		list = list[:len(list)-1]   // truncate slice
	}
	return list
}

func (this *Room) SendMessage(id int, message interface{}) {
	if client, found := this.clientsById[id]; found {
		client.SendMessage(message)
	}
}

func (this *Room) Broadcast(client *Client, message interface{}) {
	for id,c := range this.clientsById {
		if id != client.Id {
			c.SendMessage(message)
		}
	}
}

func (this *Room) BroadcastToViewers(message interface{}) {
	for _,c := range this.Clients["viewer"] {
		c.SendMessage(message)
	}
}

func (this *Room) SendMessageToEngine(message interface{}) {
	this.SendMessage(this.engine.Id, message)
}
