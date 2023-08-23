package lib

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-jwt-auth/internal/services"
	"go-jwt-auth/internal/storage"
	"io"
	"os"
)

const (
	_defaultConfigPath = "config/config.json"
)

type Config struct {
	Port  string `json:"port"`
	HTTPS bool   `json:"https"`

	PathToConfig string `json:"-"`

	StorageConfig storage.Config     `json:"storage"`
	JWTConfig     services.JWTConfig `json:"jwt"`
}

// New creates a new config.
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
