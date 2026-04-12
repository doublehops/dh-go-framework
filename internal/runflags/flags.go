// Package runflags parses and exposes CLI flags used when running the server.
package runflags

import "flag"

// FlagStruct holds the parsed CLI flag values for the server.
type FlagStruct struct {
	ConfigFile string
}

// GetFlags parses and returns the CLI flags provided at startup.
func GetFlags() FlagStruct {
	configFile := flag.String("config", "config.json", "Config file to use")
	flag.Parse()

	params := FlagStruct{
		ConfigFile: *configFile,
	}

	return params
}
