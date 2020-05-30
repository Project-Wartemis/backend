package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	log "github.com/sirupsen/logrus"
	msg "github.com/Project-Wartemis/pw-backend/internal/message"
	"github.com/Project-Wartemis/pw-backend/internal/util"
)

const (
	PLAYER_PREFIX = "{[(___"
	PLAYER_SUFFIX = "___)]}"
)

var (
	GAME_COUNTER util.SafeCounter
)

type stateConverter func(*msg.StateMessage) *msg.StateMessageOut

type Game struct {
	Room
	Id int           `json:"id"`
	Engine *Client   `json:"engine"`
	Players []*Player  `json:"players"`
	History *History `json:"-"`
	Started bool     `json:"started"`
	Stopped bool     `json:"stopped"`
}

func NewGame(name string, engine *Client) *Game {
	room := NewRoom(name)
	return &Game {
		Room: *room,
		Id: GAME_COUNTER.GetNext(),
		Engine: engine,
		Players: []*Player{},
		History: NewHistory(),
		Started: false,
		Stopped: false,
	}
}

func (this *Game) Start() error {
	if this.GetStarted() {
		return errors.New(fmt.Sprintf("Game [%d] has already started", this.GetId()))
	}
	this.setStarted(true)
	players := this.GetPlayerIds()
	message := msg.NewStartMessage(this.GetId(), players, PLAYER_PREFIX, PLAYER_SUFFIX)
	this.getEngine().SendMessage(message)
	return nil
}

func (this *Game) Stop() error {
	if !this.GetStarted() {
		return errors.New(fmt.Sprintf("Game [%d] has not started yet", this.GetId()))
	}
	if this.GetStopped() {
		return errors.New(fmt.Sprintf("Game [%d] has already stopped", this.GetId()))
	}
	this.setStopped(true)
	message := msg.NewStopMessage(this.GetId())
	this.BroadcastToType(TYPE_BOT, message)
	return nil
}



// communication related stuff

func (this *Game) HandleReconnect(client *Client) {
	// resend the last history state
	if client.GetType() != TYPE_BOT {
		return // only to bots
	}
	players := this.getPlayersByClient(client)
	if len(players) == 0 {
		return // not a player in this game
	}

	message := this.GetHistory().GetLatest()
	if message == nil {
		return // no history yet
	}

	for _,player := range players {
		log.Debugf("Sending last state in game [%s] to [%s]", this.GetName(), client.GetName())
		this.sendStateMessageToPlayer(player, message)
	}
}

func (this *Game) HandleStateMessage(message *msg.StateMessage) {
	for _,player := range this.Players {
		this.sendStateMessageToPlayer(player, message)
	}
	broadcast := this.makeStateConverter(nil)(message)
	this.BroadcastToType(TYPE_VIEWER, broadcast)
	this.GetHistory().Add(message)
	this.GetHistory().AddConverted(broadcast)
}

func (this *Game) makeStateConverter(player *Player) stateConverter {
	paddedPlayerId := this.getPaddedId(-1)
	playerKey := ""
	if player != nil {
		paddedPlayerId = this.getPaddedId(player.GetId())
		playerKey = player.GetKey()
	}
	regex1 := regexp.MustCompile(regexp.QuoteMeta("\"" + paddedPlayerId + "\""))
	regex2 := regexp.MustCompile(regexp.QuoteMeta("\"" + PLAYER_PREFIX) + "(\\d+)" + regexp.QuoteMeta(PLAYER_SUFFIX + "\""))

	return func(message *msg.StateMessage) *msg.StateMessageOut {
		state := string(message.State)
		move := !this.GetStopped() && util.Includes(message.Players, paddedPlayerId)
		state = regex1.ReplaceAllString(state, "1")
		state = regex2.ReplaceAllString(state, "$1")
		return msg.NewStateMessageOut(message.Game, playerKey, message.Turn, move, state)
	}
}

func (this *Game) sendStateMessageToPlayer(player *Player, message *msg.StateMessage) {
	outgoing := this.makeStateConverter(player)(message)
	player.GetClient().SendMessage(outgoing)
}

func (this *Game) HandleActionMessage(message *msg.ActionMessage) {
	player := this.GetPlayerByKey(message.Key)
	if player == nil {
		return
	}
	if !this.GetStarted() {
		player.GetClient().SendError(fmt.Sprintf("Game [%d] has not started yet", this.GetId()))
		return
	}
	if this.GetStopped() {
		player.GetClient().SendError(fmt.Sprintf("Game [%d] has already stopped", this.GetId()))
		return
	}
	message.Player = this.getPaddedId(player.GetId())
	this.getEngine().SendMessage(message)
}

// getters and setters

func (this *Game) GetId() int {
	this.RLock()
	defer this.RUnlock()
	return this.Id
}

func (this *Game) getPaddedId(id int) string {
	return PLAYER_PREFIX + strconv.Itoa(id) + PLAYER_SUFFIX
}

func (this *Game) getEngine() *Client {
	this.RLock()
	defer this.RUnlock()
	return this.Engine
}

func (this *Game) AddPlayer(client *Client) {
	if this.GetStarted() {
		return
	}

	log.Infof("Adding player [%s] to game [%s]", client.GetName(), this.GetName())

	player := NewPlayer(this.GetNextPlayerId(), client)
	this.Lock()
	defer this.Unlock()
	this.Players = append(this.Players, player)
}

func (this *Game) RemovePlayer(player *Player) {
	log.Infof("Removing player [%s] from game [%s]", player.GetKey(), this.GetName())

	this.Lock()
	defer this.Unlock()
	for i,p := range this.Players {
		if p.GetId() == player.GetId() {
			this.Players[i] = this.Players[len(this.Players)-1] // copy last element to index i
			this.Players[len(this.Players)-1] = nil             // erase last element
			this.Players = this.Players[:len(this.Players)-1]   // truncate slice
			return
		}
	}
}

func (this *Game) GetPlayerIds() []int {
	this.RLock()
	defer this.RUnlock()
	result := []int{}
	for _,player := range this.Players {
		result = append(result, player.GetId())
	}
	return result
}

func (this *Game) GetNextPlayerId() int {
	this.RLock()
	defer this.RUnlock()
	if len(this.Players) == 0 {
		return 2
	}
	return this.Players[len(this.Players)-1].GetId() + 1
}

func (this *Game) GetPlayerByKey(key string) *Player {
	this.RLock()
	defer this.RUnlock()
	for _,player := range this.Players {
		if player.GetKey() == key {
			return player
		}
	}
	log.Errorf("Could not find player in game [%s] with key [%s]. This is unexpected", this.GetName(), key)
	return nil
}

func (this *Game) getPlayersByClient(client *Client) []*Player {
	result := []*Player{}
	this.RLock()
	defer this.RUnlock()
	for _,player := range this.Players {
		log.Infof("comparing [%p] to [%p]", player.GetClient(), client)
		if player.GetClient() == client {
			result = append(result, player)
		}
	}
	return result
}

func (this *Game) GetHistory() *History {
	this.RLock()
	defer this.RUnlock()
	return this.History
}

func (this *Game) GetStarted() bool {
	this.RLock()
	defer this.RUnlock()
	return this.Started
}

func (this *Game) setStarted(started bool) {
	this.Lock()
	defer this.Unlock()
	this.Started = started
}

func (this *Game) GetStopped() bool {
	this.RLock()
	defer this.RUnlock()
	return this.Stopped
}

func (this *Game) setStopped(stopped bool) {
	this.Lock()
	defer this.Unlock()
	this.Stopped = stopped
}



// lock for json marshalling

type JGame Game

func (this *Game) MarshalJSON() ([]byte, error) {
    this.RLock()
    defer this.RUnlock()
    return json.Marshal(JGame(*this))
}
