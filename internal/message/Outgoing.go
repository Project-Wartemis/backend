package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type GamestateMessage struct {
	Type string                        `json:"type"`
	Payload map[string]json.RawMessage `json:"payload"`
}

type MoveRequestMessage struct {
	Type string                        `json:"type"`
	Key string                         `json:"-"`
	Payload map[string]json.RawMessage `json:"payload"`
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
