package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	_defaultConfigPath = "config.json"
)

// Flag is a flag for config.
type Flag struct {
	PathToConfig *string
}

var _f = &Flag{}

type Config struct {
	Host  string `json:"host"`
	HTTPS bool   `json:"https"`

	PathToConfig string `json:"-"`

	DSN string `json:"dsn"`
}

func init() {
	_f.PathToConfig = flag.String("config", _defaultConfigPath, "-config=path/to/conf.json")
}

// New creates a new config.
func New() (*Config, error) {
	flag.Parse()

	if _f.PathToConfig == nil {
		return nil, fmt.Errorf("config file is required")
	}

	if err := Modify(*_f.PathToConfig); err != nil {
		return nil, fmt.Errorf("can't modify config: %v", err)
	}

	return &Config{}, nil
}

// Modify modifies the config by the file provided.
func Modify(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("can't open %s: %v", filename, err)
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read %s: %v", filename, err)
	}

	var fcopy Config
	err = json.Unmarshal(all, &fcopy)
	if err != nil {
		return fmt.Errorf("can't unmarshal %s: %v", filename, err)
	}

	return nil
}
