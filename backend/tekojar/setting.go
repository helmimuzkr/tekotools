package tekojar

import (
	"encoding/json"
	"fmt"
	"os"
)

type Setting struct {
	path           string
	currentSetting *TekojarSetting
}

type TekojarSetting struct {
	Command         string           `json:"command"`
	AutoShutdown    bool             `json:"auto_shutdown"`
	ServiceSettings []ServiceSetting `json:"service_settings"`
}

type ServiceSetting struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	SkipFlag bool   `json:"skip_flag"`
	Idx      int    `json:"idx"`
	Delay    int    `json:"delay"`
}

func LoadSetting() (*Setting, error) {
	path, _ := ConcatWithExecutablePath("settings.json")
	PrintLog(SYSTEM, 0, "Setting Path", path)
	s := &Setting{path: path}
	if _, err := s.Load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Setting) Load() (*TekojarSetting, error) {
	if s.isFileExist() {
		if err := s.loadFile(); err != nil {
			return nil, err
		}
	} else {
		if err := s.createFile(); err != nil {
			return nil, err
		}
	}
	return s.currentSetting, nil
}

func (s *Setting) isFileExist() bool {
	if _, err := os.Stat(s.path); err == nil {
		return true
	} else {
		errMsg := fmt.Sprintf("is setting not exist: %v", os.IsNotExist(err))
		PrintErr(SYSTEM, 0, errMsg, err.Error())
		return false
	}
}

func (s *Setting) loadFile() error {
	f, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("failed to read setting: %w", err)
	}

	tekojarSetting := TekojarSetting{}
	if err := json.Unmarshal(f, &tekojarSetting); err != nil {
		return fmt.Errorf("failed to parse setting: %w", err)
	}

	s.currentSetting = &tekojarSetting
	return nil
}

func (s *Setting) createFile() error {
	tekojarSetting := TekojarSetting{
		ServiceSettings: []ServiceSetting{},
	}
	if err := s.Save(tekojarSetting); err != nil {
		return err
	}
	return nil
}

func (s *Setting) Save(ts TekojarSetting) error {
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal setting: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write setting: %w", err)
	}

	PrintLog(SYSTEM, 0, fmt.Sprintf("Save Setting : %#v", ts))

	s.currentSetting = &ts
	return nil
}

func (s *Setting) Current() *TekojarSetting {
	if s.currentSetting == nil {
		ts, _ := s.Load()
		s.currentSetting = ts
	}
	return s.currentSetting
}
