package util

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

func InternalServerError(writer http.ResponseWriter, message string, err error) {
	log.Errorf("%s: %s", message, err)
	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte(fmt.Sprintf(`{"message": "%s: %s"}`, message, err)))
}

func WriteJson(writer http.ResponseWriter, value interface{}) {
	json, err := json.Marshal(value)
	if err != nil {
		InternalServerError(writer, "Error parsing object to json", err)
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
