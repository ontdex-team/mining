package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	PriceApi string `json:"price_api"`
}

func ParseConfig() (*Config, error) {
	fileContent, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(fileContent, config)
	return config, err
}
