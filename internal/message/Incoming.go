package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	Type string
}

// TODO temporary, remove later
type EchoMessage struct {
	Value string
}

type RegisterMessage struct {
	Name string
	Key string
}

func ParseMessage(raw []byte) *Message {
	message := &Message{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Errorf("Could not parse Message {%s}", raw)
		log.Panic(err)
	}
	return message
}

func ParseEchoMessage(raw []byte) *EchoMessage {
	message := &EchoMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Errorf("Could not parse EchoMessage {%s}", raw)
		log.Panic(err)
	}
	return message
}

func ParseRegisterMessage(raw []byte) *RegisterMessage {
	message := &RegisterMessage{}
	err := json.Unmarshal(raw, message)
	if err != nil {
		log.Errorf("Could not parse RegisterMessage {%s}", raw)
		log.Panic(err)
	}
	return message
}
