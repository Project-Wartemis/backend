package validation

import (
	// "os"
	"io/ioutil"
	"encoding/json"
	"github.com/qri-io/jsonschema"
	"github.com/sirupsen/logrus"
)

const (
	BOT_MOVE = 0
	BOT_REGISTER = 1
	NEW_GAME_REQUEST = 2
	NEW_GAME_RESPONSE = 3
)

type Validator struct {
	schemas map[int]*jsonschema.RootSchema
}

func NewValidator() *Validator {
	v := Validator{}
	v.schemas = map[int]*jsonschema.RootSchema{
		BOT_MOVE: loadSchema("pkg/validation/resources/bots/botMove.json"),
		BOT_REGISTER: loadSchema("pkg/validation/resources/bots/botRegister.json"),
		NEW_GAME_REQUEST: loadSchema("pkg/validation/resources/RESTapi/NewGameRequest.json"),
		NEW_GAME_RESPONSE: loadSchema("pkg/validation/resources/RESTapi/NewGameResponse.json"),
	}
	return &v
}

func loadSchema(schemaPath string) *jsonschema.RootSchema{
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		logrus.Fatalf("Failed to read schema file: %s", err)
		panic(err)
	}
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		logrus.Fatalf("Unable to read schema: %s", err)
		panic(err)
	}
	return rs
}

func (v *Validator) ValidateBytes(bytes []byte, schemaEnum int) bool {
	schema := v.schemas[schemaEnum]
	mistakes, err := schema.ValidateBytes(bytes)
	if err != nil {
		logrus.Fatalf("Validation error: %s", err)
		panic(err)
	}
	logrus.Info(mistakes)
	return len(mistakes) == 0
}
