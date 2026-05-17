package tekojar

import (
	"encoding/json"
	"fmt"
	"os"
)

type Setting struct {
	Command      string           `json:"command"`
	AutoShutdown bool             `json:"auto_shutdown"`
	Services     []ServiceSetting `json:"services"`
}

type ServiceSetting struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	SkipFlag bool   `json:"skip_flag"`
}

func LoadSetting(path string) (*Setting, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read setting: %w", err)
	}

	var cfg Setting
	if err := json.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse setting: %w", err)
	}

	return &cfg, nil
}
