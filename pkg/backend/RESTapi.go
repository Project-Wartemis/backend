package communication

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/sirupsen/logrus"
)

// ------------------------------------
// GET get list of active bots
type getBotListHandler struct {
	Reception *Reception
}

func (h *getBotListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message": "Forbidden"}`))
		return
	}

	// Get list of bots
	activeBotNames := []string{}
	h.Reception.RefreshBotList()
	for botName, _ := range h.Reception.RegisteredBots {
		activeBotNames = append(activeBotNames, botName)
	}

	// Build response
	resposeBytes, err := json.Marshal(map[string][]string{"bots": activeBotNames})
	if err != nil {
		logrus.Errorf("Error listing active bots: %s", err)
		InternalServerError(w, "Error listing active bots", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resposeBytes)
}

// ------------------------------------
// POST start new game
type postNewGameHandler struct {
	Reception *Reception
}

type NewGameRequest struct {
	Bots []string
	GameEngineName string
}

func (h *postNewGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Forbidden(w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	ngr := NewGameRequest{}
	err := decoder.Decode(&ngr)
	if err != nil {
		BadRequest(w, err)
		return
	}

	err = h.Reception.StartNewGame(ngr)
	if err != nil {
		InternalServerError(w, "Could not reach GameEngine", err)
		return
	}
}

// ------------------------------------
// Start the rest API
func startRestApi(recep *Reception) {
	// List active bots
	listBots := &getBotListHandler{Reception: recep}
	listBotEndpoint := "/bot/list"
	http.Handle(listBotEndpoint, listBots)
	logrus.Infof("Listening for GET on %s", listBotEndpoint)

	// Start a new game
	newGame := &postNewGameHandler{Reception: recep}
	newGameEndpoint := "/game/new"
	http.Handle(newGameEndpoint, newGame)
	logrus.Infof("Listening for POST on %s", newGameEndpoint)
}

// ------------------------------------
// Default replies
// ------------------------------------
func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"message": "Forbidden"}`))
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf(`{"message": "BadRequest", "error": "%s"}`, err)))
	logrus.Errorf("Could not read body: %s", err)
}

func InternalServerError(w http.ResponseWriter, message string, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf(`{"message": "%s", "error": "%s"}`, message, err)))
	logrus.Errorf("%s: %s", message, err)
}
