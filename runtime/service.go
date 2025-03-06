package runtime

import (
	"context"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/di"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/pkg/errors"
	"github.com/yukinagae/genkit-go-plugins/plugins/openai"
	"gorm.io/gorm"
)

type (
	Runtime interface {
		Run(ctx context.Context, threadId uint, agentId uint) error
	}
	service struct {
		db *gorm.DB
	}
)

var (
	Key = di.NewKey()
)

func init() {
	di.Register(Key, func(c context.Context, _ *di.Container) (any, error) {
		conf, err := di.Get[*config.RuntimeConfig](c, config.RuntimeConfigKey)
		if err != nil {
			return nil, err
		}

		dbInstance, err := di.Get[*gorm.DB](c, db.Key)
		if err != nil {
			return nil, err
		}

		if err := openai.Init(c, &openai.Config{
			APIKey: conf.OpenAIApiKey,
		}); err != nil {
			return nil, errors.WithStack(err)
		}

		return &service{
			db: dbInstance,
		}, nil
	})
}
