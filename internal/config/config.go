package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Host    Host    `json:"host"`
	Logging Logging `json:"logging"`
	DB      DB      `json:"database"`
}

type Host struct {
	Port string `json:"port"`
}

type Aggregator struct {
	Name string `json:"name"`
}

type Logging struct {
	Writer       string
	LogLevel     string `json:"logLevel"`
	OutputFormat string `json:"outputFormat"`
}

type DB struct {
	User string `json:"user"`
	Pass string `json:"password"`
	Host string `json:"host"`
	Name string `json:"name"`
}

func New(configFile string) (*Config, error) {
	log.Printf("Loading config from file: %s", configFile)

	var c Config

	pwd, _ := os.Getwd()
	relPath := pwd + "/" + configFile

	f, err := os.ReadFile(relPath)
	if err != nil {
		log.Printf("unable to read config file - %s. %s", relPath, err.Error())

		return nil, fmt.Errorf("unable to read config file `%s`. %w", configFile, err)
	}

	if err := json.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	if c.DB.Host == "" || c.DB.Name == "" || c.DB.User == "" || c.DB.Pass == "" {
		return &c, fmt.Errorf("some configuration is missing")
	}

	return &c, nil
}
