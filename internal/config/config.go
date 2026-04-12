package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Host    Host    `json:"host"`
	Logging Logging `json:"logging"`
	DB      DB      `json:"database"`
}

type Host struct {
	Port    string `json:"port"`
	TestURL string `json:"testUrl"`
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

	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load .env")
	}

	var c Config

	pwd, _ := os.Getwd()
	relPath := pwd + "/" + configFile
	if _, err := os.Stat(relPath); errors.Is(err, os.ErrNotExist) {
		relPath = pwd + "/../../../" + configFile // test path.
	}

	f, err := os.ReadFile(relPath)
	if err != nil {
		log.Printf("unable to read config file - %s. %s", relPath, err.Error())

		return nil, fmt.Errorf("unable to read config file `%s`. %w", configFile, err)
	}

	if err := json.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	resolveEnvInStruct(&c)

	if c.DB.Host == "" || c.DB.Name == "" || c.DB.User == "" || c.DB.Pass == "" {
		return &c, fmt.Errorf("some configuration is missing")
	}

	return &c, nil
}

// resolveEnvInStruct will recursively look through config and replace any values with `ENV:XXX` with the
// corresponding value found in `.env`.
func resolveEnvInStruct(s interface{}) {
	val := reflect.ValueOf(s)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return
	}
	val = val.Elem()
	// typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// fieldType := typ.Field(i)

		switch field.Kind() {
		case reflect.String:
			if field.CanSet() && strings.HasPrefix(field.String(), "ENV:") {
				envKey := strings.TrimPrefix(field.String(), "ENV:")
				if envVal, ok := os.LookupEnv(envKey); ok {
					field.SetString(envVal)
				} else {
					fmt.Printf("Warning: env var %s not set\n", envKey)
				}
			}
		case reflect.Struct:
			resolveEnvInStruct(field.Addr().Interface())
		}
	}
}
