package base

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type AuthorizationConfig struct {
	AccessToken string `yaml:"accessToken" validate:"required"`
}

type ServerConfig struct {
	Socket                 string `yaml:"socket" validate:"required,unix_addr"`
	BasePath               string `yaml:"basePath"`
	OpenapiBasePath        string `yaml:"openapiBasePath"`
	PaginationDefaultLimit int64  `yaml:"paginationDefaultLimit" validate:"required,gt=1"`
}

type KratosConfig struct {
	AdminApiUrl string `yaml:"adminApiUrl" validate:"required"`
}

type LogConfig struct {
	Level   string `yaml:"level" validate:"required,oneof=fatal error warn warning info debug trace"`
	AppName string `yaml:"appName" validate:"required"`
}

type BackendConfig struct {
	Server ServerConfig        `yaml:"server"`
	Logs   LogConfig           `yaml:"logs"`
	Kratos KratosConfig        `yaml:"kratos"`
	Auth   AuthorizationConfig `yaml:"authorization"`
}

func LoadConfiguration(file string) (*BackendConfig, error) {
	Logger.WithFields(logrus.Fields{"filename": file}).Info(
		"Loading configuration",
	)

	cfg := &BackendConfig{}
	cfg.SetDefaults()
	if err := cfg.loadFromFile(file); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *BackendConfig) SetDefaults() {
	cfg.Server.Socket = "localhost:8000"
	cfg.Server.BasePath = "/backend"
	cfg.Server.OpenapiBasePath = "/swagger"
	cfg.Server.PaginationDefaultLimit = 20

	cfg.Logs.Level = logrus.DebugLevel.String()
	cfg.Logs.AppName = "sharing-backend"

	cfg.Kratos.AdminApiUrl = "http://127.0.0.1:4434"
}

func (cfg *BackendConfig) loadFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("config file '%s' open error. %s", file, err.Error())
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return fmt.Errorf(
			"config file '%s' reading error, invalid format. %s",
			file,
			err.Error(),
		)
	}
	return nil
}

func (cfg *BackendConfig) validate() error {
	validatorObj := validator.New()
	if err := validatorObj.Struct(cfg); err != nil {
		return WrapValidationErrors(err)
	}
	return nil
}
