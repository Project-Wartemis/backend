package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	Type string `json:"type"`
}

type RegisterMessage struct {
	Message
	ClientType string
	Name string
}

type RoomMessage struct {
	Message
	Name string
	Engine int
}

type InviteMessage struct { // also outgoing
	Message
	Client int  `json:"client"`
	Room int    `json:"room"`
	Name string `json:"name"`
}

type StartMessage struct { // also outgoing
	Message
	Players []int `json:"players"`
}

type StateMessage struct { // also outgoing
	Message
	Players []int                    `json:"players"`
	State map[string]json.RawMessage `json:"state"`
}

type ActionMessage struct { // also outgoing
	Message
	Player int                        `json:"player"`
	Action map[string]json.RawMessage `json:"action"`
}

func NewInviteMessage(room int, name string, client int) *InviteMessage {
	message := Message {
		Type: "invite",
	}
	return &InviteMessage {
		Message: message,
		Client: client,
		Room: room,
		Name: name,
	}
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

func ParseRegisterMessage(raw []byte) (*RegisterMessage, error) {
	message := &RegisterMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse RegisterMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseRoomMessage(raw []byte) (*RoomMessage, error) {
	message := &RoomMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse RoomMessage [%s]", raw)
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

func ParseStartMessage(raw []byte) (*StartMessage, error) {
	message := &StartMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse StartMessage [%s]", raw)
		return nil, err
	}
	return message, nil
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

func ParseActionMessage(raw []byte) (*ActionMessage, error) {
	message := &ActionMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse ActionMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}
