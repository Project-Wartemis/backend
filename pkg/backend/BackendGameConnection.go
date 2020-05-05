package backend

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/Project-Wartemis/pw-backend/pkg/validation"
)

type BotMove struct {
	AccessKey string
	Move  string
}

type BackendGameConnection struct {
	Game     *Game
	endpoint string
}

func newBackendGameConnection(game *Game, endpoint string) *BackendGameConnection {
	gc := &BackendGameConnection{
		Game:     game,
		endpoint: endpoint,
	}
	http.HandleFunc(gc.endpoint, gc.ListenBotMove)
	logrus.Infof("Opening new BackendGameConnection for endpoint %s", endpoint)
	return gc
}

// Connection is not closed
func (g *BackendGameConnection) ListenBotMove(w http.ResponseWriter, r *http.Request) {
	ws, err := GetWebsocket(w, r)
	if err != nil {
		logrus.Errorf("Could not connect: %s", err)
	}

	botMove, err := g.readAndParseBotMove(ws)
	if err != nil {
		// Could not read/parse message
		// Disqualify bot
		ws.Close()
		return
	}
	if !g.checkAccessKey(botMove.AccessKey) {
		// Wrong accessKey
		// Close connection
		ws.Close()
		return
	}
	logrus.Infof("recieved move from bot %s", botMove.AccessKey)

	//TODO send move to GameEngine
}

func (g *BackendGameConnection) readAndParseBotMove(ws *websocket.Conn) (*BotMove, error) {
	message, err := ReadMessage(ws)
	if err != nil {
		logrus.Errorf("read: %s", err)
		return nil, err
	}

	if !validator.ValidateBytes(message, validation.NEW_GAME_REQUEST) {
		return nil, errors.New("Json is not compliant with schema")
	}

	// Parse json
	botMove := BotMove{}
	err = json.Unmarshal(message, &botMove)
	if err != nil {
		logrus.Errorf("Error parsing botmove: %s", err)
		return nil, err
	}
	return &botMove, nil
}

func (gc *BackendGameConnection) checkAccessKey(accessKey string) bool {
	return gc.Game.Players[accessKey] != nil
}
