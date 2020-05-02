package communication

import (
	"github.com/sirupsen/logrus"
)

type Game struct {
	Players               map[string]*Bot
	BackendGameConnection *BackendGameConnection
	GameEngineConnection  *GameEngineConnection
	Finished              bool
	Winner                *Bot
}



func (g *Game) notifyPlayers() error {
	for _, p := range g.Players {
		err := p.SendMessage(g.BackendGameConnection.endpoint)
		if err != nil {
			logrus.Errorf("Could not notify player: %s. Error: ", p.BotName, err)
			return err
		}
	}
	return nil
}
