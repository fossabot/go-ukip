package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	KeystrokeWindow int    `json:"keystroke_window"`
	AbnormalTyping  int    `json:"abnormal_typing"`
	RunMode         string `json:"run_mode"`
}

func Load() (*Config, error) {
	file, err := os.Open("/etc/ukip/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}