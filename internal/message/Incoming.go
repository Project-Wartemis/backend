package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	Type string `json:"type"`
}

type ActionMessage struct { // also outgoing
	Message
	Game int               `json:"game"`
	Key string             `json:"key"` // incoming only
	Player string          `json:"player"` // outgoing only
	Action json.RawMessage `json:"action"`
}

type GameMessage struct {
	Message
	Name string
	Engine int
}

type InviteMessage struct {
	Message
	Game int
	Bot int
}

type JoinMessage struct {
	Message
	Game int
}

type LeaveMessage struct {
	Message
	Game int
}

type RegisterMessage struct {
	Message
	ClientType string
	Name string
	Game string
}

type StartMessage struct { // also outgoing
	Message
	Game int      `json:"game"`
	Players []int `json:"players"` // outgoing only
	Prefix string `json:"prefix"` // outgoing only
	Suffix string `json:"suffix"` // outgoing only
}

type StateMessage struct {
	Message
	Game int
	Turn int
	Players []string
	State json.RawMessage
}

type StopMessage struct { // also outgoing
	Message
	Game int `json:"game"`
}

func ParseMessage(raw []byte) (*Message, error) {
	message := &Message{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse Message [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseActionMessage(raw []byte) (*ActionMessage, error) {
	message := &ActionMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse ActionMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseGameMessage(raw []byte) (*GameMessage, error) {
	message := &GameMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse GameMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseInviteMessage(raw []byte) (*InviteMessage, error) {
	message := &InviteMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse InviteMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseJoinMessage(raw []byte) (*JoinMessage, error) {
	message := &JoinMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse JoinMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseLeaveMessage(raw []byte) (*LeaveMessage, error) {
	message := &LeaveMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse LeaveMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseRegisterMessage(raw []byte) (*RegisterMessage, error) {
	message := &RegisterMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse RegisterMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseStartMessage(raw []byte) (*StartMessage, error) {
	message := &StartMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse StartMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func NewStartMessage(game int, players []int, prefix string, suffix string) *StartMessage {
	message := Message {
		Type: "start",
	}
	return &StartMessage {
		Message: message,
		Game: game,
		Players: players,
		Prefix: prefix,
		Suffix: suffix,
	}
}

func ParseStateMessage(raw []byte) (*StateMessage, error) {
	message := &StateMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse StateMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseStopMessage(raw []byte) (*StopMessage, error) {
	message := &StopMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse StopMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func NewStopMessage(game int) *StopMessage {
	message := Message {
		Type: "stop",
	}
	return &StopMessage {
		Message: message,
		Game: game,
	}
}
