package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type GameMessage struct {
	Type string `json:"type"`
	Key string  `json:"key"`
}

type GamestateMessage struct {
	Type string                        `json:"type"`
	Payload map[string]json.RawMessage `json:"payload"`
}

type MoveRequestMessage struct {
	Type string                        `json:"type"`
	Key string                         `json:"-"`
	Payload map[string]json.RawMessage `json:"payload"`
}

func ParseGameMessage(raw []byte) *GameMessage {
	message := &GameMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Errorf("Could not parse GameMessage {%s}", raw)
		log.Panic(err)
	}
	return message
}

func ParseGamestateMessage(raw []byte) *GamestateMessage {
	message := &GamestateMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Errorf("Could not parse GamestateMessage {%s}", raw)
		log.Panic(err)
	}
	return message
}

func ParseMoveRequestMessage(raw []byte) *MoveRequestMessage {
	message := &MoveRequestMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Errorf("Could not parse MoveRequestMessage {%s}", raw)
		log.Panic(err)
	}
	return message
}
