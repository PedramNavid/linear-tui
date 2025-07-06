package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	LinearAPIKey string `json:"linear_api_key"`
	Theme        Theme  `json:"theme"`
}

type Theme struct {
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	BackgroundColor string `json:"background_color"`
	TextColor      string `json:"text_color"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "linear-tui", "config.json")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func DefaultConfig() *Config {
	return &Config{
		Theme: Theme{
			PrimaryColor:    "205",
			SecondaryColor:  "135",
			BackgroundColor: "235",
			TextColor:       "252",
		},
	}
}

func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "linear-tui")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}