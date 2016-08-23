package config

import (
	"encoding/json"
	"github.com/t11230/ramenbot/lib/modules/modulebase"
	"os"
)

type Configuration struct {
	Token   string
	MongoDB string
	Modules []modulebase.ModuleConfig
}

func LoadConfig(filename string) (*Configuration, error) {

	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)

	conf := Configuration{}
	err := decoder.Decode(&conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}
