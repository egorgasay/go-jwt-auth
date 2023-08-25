package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-jwt-auth/internal/config"
	"io"
	"os"
)

const (
	_defaultConfigPath = "config.json"
)

type Config struct {
	Port  string `json:"port"`
	HTTPS bool   `json:"https"`

	PathToConfig string `json:"-"`

	Storage config.Storage `json:"storage"`
	JWT     config.JWT     `json:"jwt"`
}

// NewConfig creates a new config.
func NewConfig() (conf Config, err error) {
	flag.Parse()

	if conf, err = FromJSON(_defaultConfigPath); err != nil {
		return conf, fmt.Errorf("can't load config: %v", err)
	}

	return conf, nil
}

// FromJSON loads a config from a json file.
func FromJSON(filename string) (conf Config, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return conf, fmt.Errorf("can't open %s: %v", filename, err)
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		return conf, fmt.Errorf("can't read %s: %v", filename, err)
	}

	err = json.Unmarshal(all, &conf)
	if err != nil {
		return conf, fmt.Errorf("can't unmarshal %s: %v", filename, err)
	}

	return conf, nil
}
