package base

import (
	"sync"
)

type Lobby struct {
	Room
	Rooms []*Room `json:"rooms"`
	roomsByKey map[string]*Room
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

	room := NewRoom("Lobby")
	lobby = &Lobby {
		Room: *room,
		Rooms: []*Room{},
		roomsByKey: map[string]*Room{},
	}

	return lobby
}

func (this *Lobby) GetRoomByKey(key string) *Room {
	return this.roomsByKey[key]
}

func (this *Lobby) AddRoom(room *Room) *Room {
	result := NewRoom(room.Name)
	this.Rooms = append(this.Rooms, result)
	this.roomsByKey[result.Key] = result
	return result
}
