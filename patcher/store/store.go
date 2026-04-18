package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type AppConfig struct {
	BasePath   string
	Submodules []string `json:"submodules"`
	Repo       string   `json:"repo"`
}

func (c *AppConfig) GetPatchesDir() string {
	return filepath.Join(c.BasePath, "patches")
}

var Config = AppConfig{}

func InitConfig(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &Config)
	if err != nil {
		return err
	}

	sort.Slice(Config.Submodules, func(i, j int) bool {
		return len(Config.Submodules[i]) > len(Config.Submodules[j])
	})

	return nil
}
