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

func ParseGameMessage(raw []byte) (*GameMessage, error) {
	message := &GameMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse GameMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseGamestateMessage(raw []byte) (*GamestateMessage, error) {
	message := &GamestateMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse GamestateMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}

func ParseMoveRequestMessage(raw []byte) (*MoveRequestMessage, error) {
	message := &MoveRequestMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse MoveRequestMessage [%s]", raw)
		return nil, err
	}
	return message, nil
}
