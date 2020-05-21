package message

type ErrorMessage struct {
	Message
	Content string `json:"message"`
}

type LobbyMessage struct {
	Message
	Lobby interface{} `json:"lobby"`
}

type ConnectedMessage struct {
	Message
}

type RegisteredMessage struct {
	Message
	Id int `json:"id"`
}

func NewErrorMessage(content string) *ErrorMessage {
	message := Message {
		Type: "error",
	}
	return &ErrorMessage {
		Message: message,
		Content: content,
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

func NewConnectedMessage() *ConnectedMessage {
	message := Message {
		Type: "connected",
	}
	return &ConnectedMessage {
		Message: message,
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
