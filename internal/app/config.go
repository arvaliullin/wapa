package app

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ComposerConfig настройки сервиса
type ComposerConfig struct {
	DbConnection      string `yaml:"database_connection"`
	Address           string `yaml:"composer_address"`
	DataPath          string `yaml:"composer_data_path"`
	NatsURL           string `yaml:"nats_url"`
	NatsSubjectRunner string `yaml:"nats_subject_runner"`
	NatsSubjectResult string `yaml:"nats_subject_result"`
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
	if dbConn, exists := os.LookupEnv("COMPOSER_DB_CONNECTION"); exists {
		config.DbConnection = dbConn
	}
	if dataPath, exists := os.LookupEnv("COMPOSER_DATA_PATH"); exists {
		config.DataPath = dataPath
	}
	if natsURL, exists := os.LookupEnv("COMPOSER_NATS_URL"); exists {
		config.NatsURL = natsURL
	}

	if natsSubjectRunner, exists := os.LookupEnv("COMPOSER_NATS_SUBJECT_RUNNER"); exists {
		config.NatsSubjectRunner = natsSubjectRunner
	}

	if natsSubjectResult, exists := os.LookupEnv("COMPOSER_NATS_SUBJECT_RESULT"); exists {
		config.NatsSubjectResult = natsSubjectResult
	}

	return nil
}
