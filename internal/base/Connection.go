package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	sync "github.com/sasha-s/go-deadlock"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type Connection struct {
	sync.RWMutex
	sendLock sync.Mutex // we cannot send two messages concurrently
	connection *websocket.Conn
	pinger *time.Ticker
}

func NewConnection(connection *websocket.Conn) *Connection {
	return &Connection {
		connection: connection,
	}
}



// basic communication related stuff

func (this *Connection) SendMessage(message interface{}) error {
	text, err := json.Marshal(message)
	if err != nil {
		return errors.New(fmt.Sprintf("Unexpected error while parsing message to json : [%s] : [%s]", err, message))
	}

	this.sendLock.Lock()
	err = this.getConnection().WriteMessage(websocket.TextMessage, text)
	this.sendLock.Unlock()

	if err != nil {
		return errors.New(fmt.Sprintf("Unexpected error while sending message : [%s] : [%s]", err, text))
	}

	return nil
}

func (this *Connection) sendPing() {
	this.sendLock.Lock()
	err := this.getConnection().WriteMessage(websocket.PingMessage, nil)
	this.sendLock.Unlock()

	if err != nil {
		log.Warnf("Unexpected error while sending ping : [%s]", err)
		return
	}
}

func (this *Connection) StartPinging() {
	pinger := time.NewTicker(30 * time.Second)
	this.setPinger(pinger)
	for {
		<- pinger.C
		this.sendPing()
	}
}

func (this *Connection) StopPinging() {
	pinger := this.getPinger()
	if pinger != nil {
		pinger.Stop()
	}
}



// getters and setters

func (this *Connection) getConnection() *websocket.Conn {
	this.RLock()
	defer this.RUnlock()
	return this.connection
}

func (this *Connection) getPinger() *time.Ticker {
	this.RLock()
	defer this.RUnlock()
	return this.pinger
}

func (this *Connection) setPinger(pinger *time.Ticker) {
	this.Lock()
	defer this.Unlock()
	this.pinger = pinger
}
