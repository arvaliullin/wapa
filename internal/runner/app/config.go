package app

import (
	"os"

	"github.com/arvaliullin/wapa/internal/app"
	"gopkg.in/yaml.v3"
)

// RunnerConfig настройки сервиса composer
type RunnerConfig struct {
	app.ServiceConfig `yaml:",inline"`
	ComposerAddress   string `yaml:"composer_address"`
	DataPath          string `yaml:"data_path"`
}

func NewRunnerConfig(configPath string) (*RunnerConfig, error) {
	config := &RunnerConfig{}
	err := config.Load(configPath)
	return config, err
}

func (config *RunnerConfig) Load(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	return config.overrideEnv()
}

func (config *RunnerConfig) overrideEnv() error {
	config.ServiceConfig.OverrideEnv()

	if composerAddress, exists := os.LookupEnv("COMPOSER_ADDRESS"); exists {
		config.ComposerAddress = composerAddress
	}

	if dataPath, exists := os.LookupEnv("DATA_PATH"); exists {
		config.DataPath = dataPath
	}

	return nil
}
