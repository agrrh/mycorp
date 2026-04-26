package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	errFailedToRead  error = errors.New("failed to read config file")
	errFailedToParse error = errors.New("failed to parse config file")
)

type Config struct {
	Tokens []string `yaml:"tokens" json:"tokens"`
}

func Load(path string) (*Config, error) {
	c := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errFailedToRead
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, errFailedToParse
	}

	return c, nil
}
