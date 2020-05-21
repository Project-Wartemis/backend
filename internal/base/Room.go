package base

import (
	"errors"
	sync "github.com/sasha-s/go-deadlock"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

var (
	ROOM_COUNTER util.SafeCounter
)

type Room struct {
	sync.RWMutex
	Id int                       `json:"id"`
	Name string                  `json:"name"`
	Clients map[string][]*Client `json:"clients"`
	Started bool                 `json:"started"`
	Stopped bool                 `json:"stopped"`
	engine *Client
	clientsById map[int]*Client
}

func NewRoom(name string) *Room {
	return &Room {
		Id: ROOM_COUNTER.GetNext(),
		Name: name,
		Clients: map[string][]*Client{},
		Started: false,
		Stopped: false,
		engine: nil,
		clientsById: map[int]*Client{},
	}
}

func (this *Room) GetClientById(id int) *Client {
	this.RLock()
	defer this.RUnlock()
	return this.clientsById[id]
}

func (this *Room) GetBotIds() []int {
	this.RLock()
	defer this.RUnlock()
	result := []int{}
	for id,client := range this.clientsById {
		if client.Type == TYPE_BOT {
			result = append(result, id)
		}
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
	this.RLock()
	if client.Type == TYPE_ENGINE && this.engine != nil {
		return errors.New("engine already registered on this room")
	}
	this.RUnlock()

	this.Lock()
	defer this.Unlock()
	this.clientsById[client.Id] = client
	if client.Type != "" {
		this.Clients[client.Type] = append(this.Clients[client.Type], client)
	}
	if client.Type == "engine" {
		this.engine = client
	}
	return nil
}

func (this *Room) RemoveClient(client *Client) {
	this.Lock()
	defer this.Unlock()
	log.Infof("Removing client [%s] from room [%s]", client.Name, this.Name)
	delete(this.clientsById, client.Id)
	if client.Type != "" {
		this.Clients[client.Type] = this.removeClientFromList(client, this.Clients[client.Type])
	}
	if client.Type == "engine" {
		this.engine = nil
	}
}

// not goroutine safe, expects caller to lock
func (this *Room) removeClientFromList(client *Client, list []*Client) []*Client {
	if client == nil {
		log.Error("Detected nil client in removeClientFromList")
		return list
	}
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
	this.RLock()
	defer this.RUnlock()
	if client, found := this.clientsById[id]; found {
		go client.SendMessage(message)
	}
}

func (this *Room) Broadcast(client *Client, message interface{}) {
	this.RLock()
	defer this.RUnlock()
	for id,c := range this.clientsById {
		if id != client.Id {
			go c.SendMessage(message)
		}
	}
}

func (this *Room) BroadcastToViewers(message interface{}) {
	this.RLock()
	defer this.RUnlock()
	for _,c := range this.Clients["viewer"] {
		go c.SendMessage(message)
	}
}

func (this *Room) SendMessageToEngine(message interface{}) {
	this.RLock()
	id := this.engine.Id
	this.RUnlock()
	this.SendMessage(id, message)
}
