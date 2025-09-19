package config

import (
	"encoding/json"
	"os"

	"gproc/pkg/types"
)

const configFile = "gproc.json"

func SaveConfig(config *types.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

func LoadConfig() (*types.Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &types.Config{
			Processes: []types.Process{},
			LogDir:    "./logs",
		}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config types.Config
	err = json.Unmarshal(data, &config)
	return &config, err
}