// Package config provides application configuration loading and environment variable resolution.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the full application configuration loaded from a JSON file.
type Config struct {
	Host    Host    `json:"host"`
	Logging Logging `json:"logging"`
	DB      DB      `json:"database"`
}

// Host holds server host configuration.
type Host struct {
	Port    string `json:"port"`
	TestURL string `json:"testUrl"`
	// CORSOrigin is the exact frontend origin allowed to call the API from a browser (e.g.
	// "https://example.com"). Empty disables the CORS middleware.
	CORSOrigin string `json:"corsOrigin"`
	// TrustForwardedFor makes IP-based controls (the login rate limiter) use the last
	// X-Forwarded-For entry. Only enable when every request arrives through one reverse proxy
	// that appends to the header (e.g. an AWS ALB).
	TrustForwardedFor bool `json:"trustForwardedFor"`
}

// Aggregator holds aggregator configuration.
type Aggregator struct {
	Name string `json:"name"`
}

// Logging holds log output configuration.
type Logging struct {
	// Writer is "stdout" (default when empty) or "file". "file" writes to FilePath and to stdout.
	Writer string `json:"writer"`
	// FilePath is the log file used when Writer is "file". Its directory is created if missing.
	FilePath     string `json:"filePath"`
	LogLevel     string `json:"logLevel"`
	OutputFormat string `json:"outputFormat"`
}

// DB holds database connection configuration.
type DB struct {
	User string `json:"user"`
	Pass string `json:"password"`
	Host string `json:"host"`
	Port string `json:"port"`
	Name string `json:"name"`
}

// defaultDBPort is used when Port is left unset.
const defaultDBPort = "3306"

// PortOrDefault returns Port, falling back to defaultDBPort when unset. Every DSN built from this
// config must go through this rather than hardcoding "3306", so a non-default port is honoured.
func (d DB) PortOrDefault() string {
	if d.Port == "" {
		return defaultDBPort
	}

	return d.Port
}

// New loads and returns a Config from the given JSON file path.
func New(configFile string) (*Config, error) {
	log.Printf("Loading config from file: %s", configFile)

	if err := loadEnv(); err != nil {
		log.Println("Unable to load .env")
	}

	var c Config

	var relPath string
	if filepath.IsAbs(configFile) {
		relPath = configFile
	} else {
		pwd, _ := os.Getwd()
		relPath = pwd + "/" + configFile
		if _, err := os.Stat(relPath); errors.Is(err, os.ErrNotExist) {
			relPath = pwd + "/../../../" + configFile // test path.
		}
	}

	f, err := os.ReadFile(relPath) //nolint:gosec
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

// loadEnv finds and loads the .env file by walking up from the current working directory.
func loadEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf(".env not found")
		}

		dir = parent
	}
}

// resolveEnvInStruct will recursively look through config and replace any values with `ENV:XXX` with the
// corresponding value found in `.env`.
func resolveEnvInStruct(s interface{}) {
	val := reflect.ValueOf(s)

	if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
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
					// Clear the placeholder so an unset optional value (e.g. CORS_ORIGIN) reads as
					// empty rather than the literal "ENV:..." string.
					log.Printf("Warning: env var %s not set\n", envKey)
					field.SetString("")
				}
			}
		case reflect.Struct:
			resolveEnvInStruct(field.Addr().Interface())
		}
	}
}
