package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	LinearAPIKey string `json:"linear_api_key"`
	Theme        Theme  `json:"theme"`
	DebugMode    bool   `json:"debug_mode"`
}

type Theme struct {
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
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
			config := DefaultConfig()
			if apiKey := os.Getenv("LINEAR_API_KEY"); apiKey != "" {
				config.LinearAPIKey = apiKey
			}
			if debugMode := os.Getenv("DEBUG"); debugMode != "" {
				config.DebugMode = true
			}
			return config, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Environment variable takes precedence
	if apiKey := os.Getenv("LINEAR_API_KEY"); apiKey != "" {
		config.LinearAPIKey = apiKey
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
