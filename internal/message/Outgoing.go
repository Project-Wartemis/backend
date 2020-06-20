package message

import (
	"encoding/json"
)

type ConnectedMessage struct {
	Message
}

type CreatedMessage struct {
	Message
	Game int `json:"game"`
}

type ErrorMessage struct {
	Message
	Error string `json:"message"`
}

type LobbyMessage struct {
	Message
	Lobby interface{} `json:"lobby"`
}

type RegisteredMessage struct {
	Message
	Id int `json:"id"`
}

type StateMessageOut struct {
	Message
	Game int              `json:"game"`
	Key string            `json:"key"`
	Turn int              `json:"turn"`
	Move bool             `json:"move"`
	State json.RawMessage `json:"state"`
}

type HistoryMessage struct {
	Message
	Messages []*StateMessageOut `json:"messages"`
}

func NewConnectedMessage() *ConnectedMessage {
	message := Message {
		Type: "connected",
	}
	return &ConnectedMessage {
		Message: message,
	}
}

func NewCreatedMessage(game int) *CreatedMessage {
	message := Message {
		Type: "created",
	}
	return &CreatedMessage {
		Message: message,
		Game: game,
	}
}

func NewErrorMessage(error string) *ErrorMessage {
	message := Message {
		Type: "error",
	}
	return &ErrorMessage {
		Message: message,
		Error: error,
	}
}

func NewLobbyMessage(lobby interface{}) *LobbyMessage {
	message := Message {
		Type: "lobby",
	}
	return &LobbyMessage {
		Message: message,
		Lobby: lobby,
	}
}

func NewRegisteredMessage(id int) *RegisteredMessage {
	message := Message {
		Type: "registered",
	}
	return &RegisteredMessage {
		Message: message,
		Id: id,
	}
}

func NewStateMessageOut(game int, key string, turn int, move bool, state string) *StateMessageOut {
	message := Message {
		Type: "state",
	}
	return &StateMessageOut {
		Message: message,
		Game: game,
		Key: key,
		Turn: turn,
		Move: move,
		State: json.RawMessage(state),
	}
}

func NewHistoryMessage(messages []*StateMessageOut) *HistoryMessage {
	message := Message {
		Type: "history",
	}
	return &HistoryMessage {
		Message: message,
		Messages: messages,
	}
}
