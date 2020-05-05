package backend

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type GameEngineConnection struct {
	GameEngine *GameEngine
	endpoint   string
}

type NewGameRegistrationForm struct {
	playerNames []string
}

func FillInNewGameRegistrationForm(players []*Bot) []byte {
	playerIds := []string{}
	for _, p := range players {
		playerIds = append(playerIds, p.AccessKey)
	}

	rf := NewGameRegistrationForm{
		playerNames: playerIds,
	}

	formBytes, err := json.Marshal(rf)
	if err != nil {
		logrus.Errorf("Failed to marshel registration form: %s", err)
	}

	return formBytes
}

func NewGameEngineConnection(ge *GameEngine, form []byte) (*GameEngineConnection, error) {
	// Make request for sending form to Game Engine
	url := "http://" + ge.address + "/register"
	logrus.Infof("Register new game at %s", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(form))
	if err != nil {
		logrus.Error("Failed create POST request for gameEngine")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Failed to send POST to gameEngine")
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	if resp.StatusCode != http.StatusAccepted {
		logrus.Errorf("Received status code %d from GameEngine", resp.Status)
		return nil, errors.New(fmt.Sprintf("Received status code %d from GameEngine", resp.Status))
	}

	type gameEngineResponse struct {
		endpoint string
	}
	ger := gameEngineResponse{}
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, ger)

	// Create the GameEngineConnection
	gec := GameEngineConnection{
		GameEngine: ge,
		endpoint: ger.endpoint,
	}
	return &gec, nil
}
