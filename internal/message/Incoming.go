package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	Type string `json:"type"`
}

// TODO temporary, remove later
type EchoMessage struct {
	Message
	Value string `json:"value"`
}

type RegisterMessage struct {
	Message
	ClientType string
	Name string
	Key string
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

func ParseEchoMessage(raw []byte) (*EchoMessage, error) {
	message := &EchoMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Warnf("Could not parse EchoMessage [%s]", raw)
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
