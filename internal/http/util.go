package http

import (
	"fmt"
	"encoding/json"
	"net/http"
	log "github.com/sirupsen/logrus"
)

// not actually used anywhere anymore
func WriteStatus(writer http.ResponseWriter, status int, message string, errors ...error) {
	log.Warnf("%s: %s", message, errors)
	writer.WriteHeader(status)
	writer.Write([]byte(fmt.Sprintf(`{"message": "%s", "errors": %s}`, message, errors)))
}

// not actually used anywhere anymore
func WriteJson(writer http.ResponseWriter, value interface{}) {
	json, err := json.Marshal(value)
	if err != nil {
		WriteStatus(writer, http.StatusBadRequest, "Error parsing object to json", err)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(json)
}
