package communication

import (
	"github.com/sirupsen/logrus"
	"time"
	"fmt"
)

type GameEngine struct {
	address string
}

func (ge *GameEngine) CreateNewGame(players []*Bot) (*Game, error) {
	// Initialise Game struct
	g := &Game{
		Players:               map[string]*Bot{},
		BackendGameConnection: nil,
		GameEngineConnection:  nil,
		Finished:              false,
		Winner:                nil,
	}
	// Add players
	for _, p := range players {
		g.Players[p.AccessKey] = p
	}
	// Create BackendGameConnection
	newEndpoint := "/" + fmt.Sprintf("%X", time.Now().UnixNano())[10:]
	g.BackendGameConnection = newBackendGameConnection(g, newEndpoint)

	// Create GameEngineConnection and register the game
	form := FillInNewGameRegistrationForm(players)
	gec, err := ge.requestConnectionWithEngine(form)
	if err != nil{
		return nil, err
	}
	g.GameEngineConnection = gec


	// Notify players
	err = g.notifyPlayers()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (ge *GameEngine) requestConnectionWithEngine(formBytes []byte) (*GameEngineConnection, error) {
	gec, err := NewGameEngineConnection(ge, formBytes)
	return gec, err
}

// For now, GameEngine definitions will be hardcoded.
// These will be read from config files in a later stage

func LoadGameEngines() map[string]*GameEngine {
	return map[string]*GameEngine{
		"demoEngine": createDemoGameEngine(),
	}
}

func createDemoGameEngine() *GameEngine{
	logrus.Info("Create demo GameEngine")
	return &GameEngine{
		address: "warlight-basic:8080",
	}
}
