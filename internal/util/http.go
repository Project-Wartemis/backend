package util

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

func WriteStatus(writer http.ResponseWriter, status int, message string, errors ...error) {
	log.Warnf("%s: %s", message, errors)
	writer.WriteHeader(status)
	writer.Write([]byte(fmt.Sprintf(`{"message": "%s", "errors": [%s]}`, message, errors)))
}

func WriteJson(writer http.ResponseWriter, value interface{}) {
	json, err := json.Marshal(value)
	if err != nil {
		WriteStatus(writer, http.StatusBadRequest, "Error parsing object to json", err)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(json)
}

func SetupWebSocket(writer http.ResponseWriter, request *http.Request, postInit func(*websocket.Conn), messageHandler func([]byte)) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Errorf("Cannot upgrade to websocket: %s", err)
		return
	}
	defer connection.Close()

	postInit(connection)

	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Errorf("Unable to read message: %s", err)
			return
		}
		messageHandler(message)
	}
}
