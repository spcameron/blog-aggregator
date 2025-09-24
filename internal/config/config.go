package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	URL             string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	filepath := filepath.Join(homedir, configFileName)

	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
