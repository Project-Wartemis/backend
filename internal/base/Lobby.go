package base

import (
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
)

func GetLobby() *Lobby {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	if lobby != nil {
		return lobby
	}

	room := NewRoom("lobby")
	lobby = &Lobby {
		Room: *room,
		Rooms: []*Room{},
		roomsById: map[int]*Room{},
	}

	return lobby
}

func (this *Lobby) GetRoomById(id int) *Room {
	return this.roomsById[id]
}

func (this *Lobby) CreateAndAddRoom(name string) *Room {
	room := NewRoom(name)
	this.AddRoom(room)
	log.Infof("Added a new room [%s]", room.Name)
	return room
}

func (this *Lobby) AddRoom(room *Room) {
	this.Rooms = append(this.Rooms, room)
	this.roomsById[room.Id] = room
	this.TriggerUpdated()
}

func (this *Lobby) RemoveRoom(room *Room) {
	delete(this.roomsById, room.Id)
	for i,r := range this.Rooms {
		if r.Id != room.Id {
			continue
		}
		this.Rooms[i] = this.Rooms[len(this.Rooms)-1] // copy last element to index i
		this.Rooms[len(this.Rooms)-1] = nil     // erase last element
		this.Rooms = this.Rooms[:len(this.Rooms)-1]   // truncate slice
	}
	this.TriggerUpdated()
}

func (this *Lobby) TriggerUpdated() {
	this.BroadcastToViewers(message.NewLobbyMessage(this))
}
