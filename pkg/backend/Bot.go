package backend

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	BotName   string
	AccessKey string
	ws        *websocket.Conn
}

func NewBot(botName string, accessKey string, ws *websocket.Conn) *Bot {
	resultBot := Bot{
		BotName:   botName,
		AccessKey: accessKey,
		ws:        ws,
	}
	if !resultBot.ping() {
		logrus.Errorf("During creation of bot, the websocket connection was closed: %s", botName)
		return nil
	}
	return &resultBot
}

func (b *Bot) ping() bool {
	return b.ws.WriteMessage(websocket.PingMessage, []byte{}) == nil
}

func (b *Bot) SendMessage(message string) error {
	err := b.ws.WriteMessage(
		websocket.TextMessage,
		[]byte(message))
	return err
}

func (b *Bot) destroy() {
	b.ws.Close()
}
