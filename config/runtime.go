package config

import (
	"context"
	goconfig "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/habiliai/agentruntime/di"
	"github.com/pkg/errors"
	"os"
)

type RuntimeConfig struct {
	LogLevel     string `env:"LOG_LEVEL"`
	LogHandler   string `env:"LOG_HANDLER"`
	OpenAIApiKey string `env:"OPENAI_API_KEY"`
	DatabaseUrl  string `env:"DATABASE_URL"`
}

var (
	RuntimeConfigKey = di.NewKey()
)

func resolveRuntimeConfig(testing bool) (*RuntimeConfig, error) {
	configReader := goconfig.New()
	if err := configReader.Feed(); err != nil {
		return nil, errors.Wrapf(err, "failed to load config")
	}

	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		configReader.AddFeeder(feeder.DotEnv{Path: ".env"})
	}
	if testing {
		filename := ".env.test"
		if v := os.Getenv("ENV_TEST_FILE"); v != "" {
			filename = v
		}
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "should be existed %s in testing", filename)
		}
		configReader.AddFeeder(feeder.DotEnv{Path: filename})
	}
	configReader.AddFeeder(feeder.Env{})

	c := RuntimeConfig{
		DatabaseUrl: "postgres://postgres:postgres@localhost:5432/test?search_path=agentruntime",
		LogLevel:    "debug",
		LogHandler:  "default",
	}
	if err := configReader.AddStruct(&c).Feed(); err != nil {
		return nil, err
	}

	return &c, nil
}

func init() {
	di.Register(RuntimeConfigKey, func(ctx context.Context, c *di.Container) (any, error) {
		return resolveRuntimeConfig(c.Env == di.EnvTest)
	})
}
