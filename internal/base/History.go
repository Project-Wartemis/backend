package base

import (
	"fmt"
	sync "github.com/sasha-s/go-deadlock"
	msg "github.com/Project-Wartemis/pw-backend/internal/message"
)

type History struct {
	sync.RWMutex
	messages []*msg.StateMessage
}

func NewHistory() *History {
	return &History {
		messages: []*msg.StateMessage{},
	}
}



// communication related stuff

func (this *History) SendAllToClient(client *Client) {
	this.RLock()
	defer this.RUnlock()
	if len(this.messages) == 0 {
		return
	}
	client.SendMessage(msg.NewHistoryMessage(this.messages))
}

func (this *History) SendTurnToClient(client *Client, turn int) {
	this.RLock()
	defer this.RUnlock()
	if turn >= len(this.messages) {
		client.SendError(fmt.Sprintf("Turn %d is not available", turn))
		return
	}
	client.SendMessage(this.messages[turn])
}

func (this *History) SendTurnsToClient(client *Client, from int, to int) {
	this.RLock()
	defer this.RUnlock()
	from = max(0, from)
	to = min(to+1, len(this.messages))
	if to <= from {
		return
	}
	slice := this.messages[from:to]
	client.SendMessage(msg.NewHistoryMessage(slice))
}



// getters and setters

func (this *History) Add(message *msg.StateMessage) {
	this.Lock()
	defer this.Unlock()
	if message.Turn >= cap(this.messages) {
		newCapacity := max(message.Turn+1, 2*cap(this.messages)) // atleast double, and more if needed
		newMessages := make([]*msg.StateMessage, len(this.messages), newCapacity)
		copy(newMessages, this.messages)
		this.messages = newMessages
	}
	this.messages = this.messages[:message.Turn+1]
	this.messages[message.Turn] = message
}



// util

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
