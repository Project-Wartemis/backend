package base

import (
	"encoding/json"
	"sync"
	log "github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/internal/message"
)

type Lobby struct {
	Room
	Rooms []*Room `json:"rooms"`
	roomsById map[int]*Room
}

var (
	lobby *Lobby
	once sync.Once
)

func initialiseLobby() {
	room := NewRoom("lobby", true)
	lobby = &Lobby{
		Room: *room,
		Rooms: []*Room{},
		roomsById: map[int]*Room{},
	}
}

func GetLobby() *Lobby {
	once.Do(initialiseLobby)
	return lobby
}



// basic communication related stuff

func (this *Lobby) TriggerUpdated() {
	this.BroadcastToType(TYPE_VIEWER, message.NewLobbyMessage(this))
}



// getters and setters

func (this *Lobby) AddRoom(room *Room) {
	this.SetRoomById(room.GetId(), room)

	this.Lock()
	defer this.Unlock()
	this.Rooms = append(this.Rooms, room)

	go this.TriggerUpdated()
}

func (this *Lobby) RemoveRoom(room *Room) {
	log.Infof("Removing room [%s] from [%s]", room.GetName(), this.GetName())

	this.RemoveRoomById(room.GetId())

	this.Lock()
	defer this.Unlock()
	for i,r := range this.Rooms {
		if r.Id == room.Id {
			this.Rooms[i] = this.Rooms[len(this.Rooms)-1] // copy last element to index i
			this.Rooms[len(this.Rooms)-1] = nil           // erase last element
			this.Rooms = this.Rooms[:len(this.Rooms)-1]   // truncate slice
		}
	}

	go this.TriggerUpdated()
}

func (this *Lobby) GetRoomById(id int) *Room {
	this.RLock()
	defer this.RUnlock()
	return this.roomsById[id]
}

func (this *Lobby) SetRoomById(id int, client *Room) {
	this.Lock()
	defer this.Unlock()
	this.roomsById[id] = client
}

func (this *Lobby) RemoveRoomById(id int) {
	this.Lock()
	defer this.Unlock()
	delete(this.roomsById, id)
}



// lock for json marshalling

type JLobby Lobby

func (this *Lobby) MarshalJSON() ([]byte, error) {
    this.RLock()
    defer this.RUnlock()
    return json.Marshal(JLobby(*this))
}
