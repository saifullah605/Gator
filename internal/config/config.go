package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, configFileName), nil

}

type Config struct {
	DBURL        string `json:"db_url"`
	CurrUserName string `json:"current_user_name"`
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(fullPath, jsonData, 0644); err != nil {
		return err
	}

	return nil

}

func Read() (Config, error) {
	fullFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(fullFilePath)
	if err != nil {
		return Config{}, err
	}

	var configs Config

	if err := json.Unmarshal(data, &configs); err != nil {
		return Config{}, err
	}

	return configs, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrUserName = name
	return write(*cfg)
}
