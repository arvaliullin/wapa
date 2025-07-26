package app

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceConfig.
type ServiceConfig struct {
	DbConnection      string `yaml:"database_connection"`
	NatsURL           string `yaml:"nats_url"`
	NatsSubjectRunner string `yaml:"nats_subject_runner"`
	NatsSubjectResult string `yaml:"nats_subject_result"`
}

func NewServiceConfig(configPath string) (*ServiceConfig, error) {
	config := &ServiceConfig{}
	err := config.Load(configPath)
	return config, err
}

func (config *ServiceConfig) Load(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	return config.OverrideEnv()
}

func (config *ServiceConfig) OverrideEnv() error {
	if dbConn, exists := os.LookupEnv("PSQL_DB_CONNECTION"); exists {
		config.DbConnection = dbConn
	}
	if natsURL, exists := os.LookupEnv("NATS_URL"); exists {
		config.NatsURL = natsURL
	}
	if natsSubjectRunner, exists := os.LookupEnv("NATS_SUBJECT_RUNNER"); exists {
		config.NatsSubjectRunner = natsSubjectRunner
	}

	if natsSubjectResult, exists := os.LookupEnv("NATS_SUBJECT_RESULT"); exists {
		config.NatsSubjectResult = natsSubjectResult
	}
	return nil
}
