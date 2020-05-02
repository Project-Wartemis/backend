package communication

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

func GetWebsocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Cannot upgrade to websocket: %s", err)
		return nil, err
	}
	return c, nil
}

func ReadMessage(ws *websocket.Conn) ([]byte, error) {
	mtype, message, err := ws.ReadMessage()
	// Check if error occured during read
	if err != nil {
		return nil, err
	}
	// Check if message is of correct type
	if mtype != websocket.TextMessage {
		logrus.Errorf("Read error: Expected textMessage but got binary message", err)
		return nil, err
	}
	return message, nil
}
