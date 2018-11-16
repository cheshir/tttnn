package main

import (
	"github.com/cheshir/tttnn/server/ai"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	envPrefix = "APP"
)

type appConfig struct {
	Host      string
	Port      string
	StaticDir string `envconfig:"static_dir"`
	AI        ai.Config
}

func loadConfig() (appConfig, error) {
	config := appConfig{}
	err := envconfig.Process(envPrefix, &config)

	return config, errors.Wrap(err, "failed to populate config")
}
