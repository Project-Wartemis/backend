package base

import (
	"fmt"
	sync "github.com/sasha-s/go-deadlock"
	msg "github.com/Project-Wartemis/pw-backend/internal/message"
)

type History struct {
	sync.RWMutex
	messages []*msg.StateMessage
	messagesConverted []*msg.StateMessageOut
}

func NewHistory() *History {
	return &History {
		messages: []*msg.StateMessage{},
		messagesConverted: []*msg.StateMessageOut{},
	}
}



// communication related stuff

func (this *History) SendAllToViewer(client *Client) {
	this.RLock()
	defer this.RUnlock()
	if len(this.messagesConverted) == 0 {
		return
	}
	client.SendMessage(msg.NewHistoryMessage(this.messagesConverted))
}

// unused
func (this *History) SendTurnsToViewer(client *Client, from int, to int) {
	this.RLock()
	defer this.RUnlock()
	from = max(0, from)
	to = min(to+1, len(this.messagesConverted))
	if to <= from {
		return
	}
	slice := this.messagesConverted[from:to]
	client.SendMessage(msg.NewHistoryMessage(slice))
}

// unused
func (this *History) SendTurnToViewer(client *Client, turn int) {
	this.RLock()
	defer this.RUnlock()
	if turn >= len(this.messagesConverted) {
		client.SendError(fmt.Sprintf("Turn %d is not available", turn))
		return
	}
	client.SendMessage(this.messagesConverted[turn])
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

func (this *History) AddConverted(message *msg.StateMessageOut) {
	this.Lock()
	defer this.Unlock()
	if message.Turn >= cap(this.messagesConverted) {
		newCapacity := max(message.Turn+1, 2*cap(this.messagesConverted)) // atleast double, and more if needed
		newMessages := make([]*msg.StateMessageOut, len(this.messagesConverted), newCapacity)
		copy(newMessages, this.messagesConverted)
		this.messagesConverted = newMessages
	}
	this.messagesConverted = this.messagesConverted[:message.Turn+1]
	this.messagesConverted[message.Turn] = message
}

func (this *History) GetLatest() *msg.StateMessage {
	this.RLock()
	defer this.RUnlock()
	if len(this.messages) == 0 {
		return nil
	}
	return this.messages[len(this.messages)-1]
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
