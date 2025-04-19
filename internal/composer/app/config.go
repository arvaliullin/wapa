package app

import (
	"os"

	"github.com/arvaliullin/wapa/internal/app"
	"gopkg.in/yaml.v3"
)

// ComposerConfig настройки сервиса composer
type ComposerConfig struct {
	app.ServiceConfig `yaml:",inline"`
	Address           string `yaml:"composer_address"`
	DataPath          string `yaml:"composer_data_path"`
}

func NewComposerConfig(configPath string) (*ComposerConfig, error) {
	config := &ComposerConfig{}
	err := config.Load(configPath)
	return config, err
}

func (config *ComposerConfig) Load(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	return config.overrideEnv()
}

func (config *ComposerConfig) overrideEnv() error {
	config.ServiceConfig.OverrideEnv()
	if dataPath, exists := os.LookupEnv("COMPOSER_DATA_PATH"); exists {
		config.DataPath = dataPath
	}
	return nil
}
