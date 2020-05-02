package communication

import (
	"encoding/json"
	"net/http"
	"errors"
	"fmt"

	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Reception struct {
	RegisteredBots map[string]*Bot
	GameEngines map[string]*GameEngine
	games       []*Game
}

func newReception() *Reception {
	r := &Reception{
		RegisteredBots: map[string]*Bot{},
	}
	r.GameEngines = LoadGameEngines()
	http.HandleFunc("/register", r.ListenBotRegister)
	logrus.Info("Listening at /register for new bots")
	return r
}


func (recep *Reception) RefreshBotList() {
	for name, bot := range recep.RegisteredBots {
		if !bot.ping() {
			logrus.Debugf("Could not reach bot: %s", name)
			bot.destroy()
			delete(recep.RegisteredBots, name)
		}
	}
}

// Connection is not closed by function
func (recep *Reception) ListenBotRegister(w http.ResponseWriter, r *http.Request) {
	ws, err := GetWebsocket(w, r)
	if err != nil {
		logrus.Errorf("Could not connect: %s", err)
	}

	bot, err := recep.ReadAndParseRegistration(ws)
	if err != nil {
		// Could not read/parse message
		logrus.Errorf("Failed to register bot: %s", err)
		_ = bot.SendMessage("Failed to parse registration request")
		return
	}

	// Check if bot already exists
	if recep.RegisteredBots[bot.BotName] != nil {
		logrus.Infof(
			"Atempt to register bot Failed: Botname %s already registered.",
			bot.BotName)
		_ = bot.SendMessage("Botname already taken")
		return
	}

	// Generate access key
	bot.AccessKey = guuid.New().String()
	responseToBot := map[string]string{
		"accessKey": bot.AccessKey,
	}
	responseBytes, err := json.Marshal(responseToBot)
	err = bot.SendMessage(string(responseBytes))
	if err != nil {
		logrus.Errorf("Failed sending ID to bot: %s", bot.BotName)
		return
	}

	recep.RegisteredBots[bot.BotName] = bot
	logrus.Infof("Bot succesfully registered: %s (%s)", bot.BotName, bot.AccessKey)
	logrus.Debugf("Number of registered bots: %d", len(recep.RegisteredBots))
}

func (recep *Reception) ReadAndParseRegistration(ws *websocket.Conn) (*Bot, error) {
	message, err := ReadMessage(ws)
	// Check if error occured during read
	if err != nil {
		logrus.Errorf("Error reading websocket: %s", err)
		return nil, err
	}

	botBuffer := Bot{
		ws: ws,
	}
	// Parse json
	err = json.Unmarshal(message, &botBuffer)
	if err != nil {
		logrus.Errorf("JsonParseException: %s", err)
		return nil, err
	}

	resultBot := NewBot(botBuffer.BotName, botBuffer.AccessKey, ws)

	return resultBot, nil
}

func (recep *Reception) StartNewGame(ngr NewGameRequest) error {
	recep.RefreshBotList()
	participatingBots := []*Bot{}
	for _, botName := range ngr.Bots {
		if recep.RegisteredBots[botName] == nil {
			logrus.Infof("Participating bot no longer available: %s", botName)
			return errors.New( fmt.Sprintf("Participating bot no longer available: %s", botName) )
		}
		participatingBots = append(
			participatingBots, recep.RegisteredBots[botName])
	}

	newGame, err := recep.GameEngines[ngr.GameEngineName].CreateNewGame(
		participatingBots )
	if err != nil {
		logrus.Errorf("Failed to create new game: %s", err)
		return err
	}
	recep.games = append(recep.games, newGame)
	return nil
}
