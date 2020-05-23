package base

import (
	"encoding/json"
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
	Id int            `json:"id"`
	Name string       `json:"name"`
	Clients []*Client `json:"clients"`
	History *History  `json:"-"`
	Started bool      `json:"started"`
	Stopped bool      `json:"stopped"`
	isLobby bool
	engine *Client
	clientsById map[int]*Client
}

func NewRoom(name string, isLobby bool) *Room {
	return &Room {
		Id: ROOM_COUNTER.GetNext(),
		Name: name,
		Clients: []*Client{},
		History: NewHistory(),
		Started: false,
		Stopped: false,
		isLobby: isLobby,
		engine: nil,
		clientsById: map[int]*Client{},
	}
}

func (this *Room) CreateAndAddClient() *Client {
	client := NewClient(this)
	this.AddClient(client)
	log.Infof("Added a new client to room [%s]", this.GetName())
	return client
}



// basic communication related stuff

func (this *Room) SendMessage(id int, message interface{}) {
	client := this.GetClientById(id)
	if client != nil {
		client.SendMessage(message)
	}
}

func (this *Room) Broadcast(id int, message interface{}) {
	this.RLock()
	defer this.RUnlock()
	for _,client := range this.Clients {
		if client.GetId() != id {
			client.SendMessage(message)
		}
	}
}

func (this *Room) BroadcastToType(Type string, message interface{}) {
	clients := this.GetClientIdsByType(Type)
	for _,id := range clients {
		this.SendMessage(id, message)
	}
}



// getters and setters

func (this *Room) GetId() int {
	this.RLock()
	defer this.RUnlock()
	return this.Id
}

func (this *Room) SetId(id int) {
	this.Lock()
	defer this.Unlock()
	this.Id = id
}

func (this *Room) GetName() string {
	this.RLock()
	defer this.RUnlock()
	return this.Name
}

func (this *Room) SetName(name string) {
	this.Lock()
	defer this.Unlock()
	this.Name = name
}

func (this *Room) GetClientIdsByType(Type string) []int {
	this.RLock()
	defer this.RUnlock()
	result := []int{}
	for id,client := range this.clientsById {
		if client.Type == Type {
			result = append(result, id)
		}
	}
	return result
}

func (this *Room) AddClient(client *Client) error {
	if client.GetType() == TYPE_ENGINE && this.getEngine() != nil {
		return errors.New("engine already registered on this room")
	}

	this.SetClientById(client.GetId(), client)
	if client.GetType() == TYPE_ENGINE && !this.GetIsLobby() {
		this.setEngine(client)
	}

	this.Lock()
	defer this.Unlock()
	this.Clients = append(this.Clients, client)

	return nil
}

func (this *Room) RemoveClient(client *Client) {
	log.Infof("Removing client [%s] from room [%s]", client.GetName(), this.GetName())

	this.RemoveClientById(client.GetId())
	if client.GetType() == TYPE_ENGINE {
		this.setEngine(nil)
	}

	this.Lock()
	defer this.Unlock()
	for i,c := range this.Clients {
		if c.Id == client.Id {
			this.Clients[i] = this.Clients[len(this.Clients)-1] // copy last element to index i
			this.Clients[len(this.Clients)-1] = nil             // erase last element
			this.Clients = this.Clients[:len(this.Clients)-1]   // truncate slice
			return;
		}
	}
}

func (this *Room) GetStarted() bool {
	this.RLock()
	defer this.RUnlock()
	return this.Started
}

func (this *Room) SetStarted(started bool) {
	this.Lock()
	defer this.Unlock()
	this.Started = started
}

func (this *Room) GetStopped() bool {
	this.RLock()
	defer this.RUnlock()
	return this.Stopped
}

func (this *Room) SetStopped(stopped bool) {
	this.Lock()
	defer this.Unlock()
	this.Stopped = stopped
}

func (this *Room) GetIsLobby() bool {
	this.RLock()
	defer this.RUnlock()
	return this.Stopped
}

// SetIsLobby not implemented, it should not update

func (this *Room) getEngine() *Client {
	this.RLock()
	defer this.RUnlock()
	return this.engine
}

func (this *Room) setEngine(engine *Client) {
	this.Lock()
	defer this.Unlock()
	this.engine = engine
}

func (this *Room) GetClientById(id int) *Client {
	this.RLock()
	defer this.RUnlock()
	return this.clientsById[id]
}

func (this *Room) SetClientById(id int, client *Client) {
	this.Lock()
	defer this.Unlock()
	this.clientsById[id] = client
}

func (this *Room) RemoveClientById(id int) {
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
