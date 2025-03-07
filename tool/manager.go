package tool

import (
	"context"
	"github.com/firebase/genkit/go/ai"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/di"
	"github.com/habiliai/agentruntime/entity"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/habiliai/agentruntime/internal/mylog"
	"gorm.io/gorm"
)

type (
	Manager interface {
		GetLocalTool(ctx context.Context, toolName string) ai.Tool
		InitializeTools(ctx context.Context) error
		GetTools(ctx context.Context, names []string) ([]entity.Tool, error)

		GetWeather(ctx context.Context, req *GetWeatherRequest) (*GetWeatherResponse, error)
	}
)

var (
	ManagerKey = di.NewKey()
)

func init() {
	di.Register(ManagerKey, func(ctx context.Context, c *di.Container) (any, error) {
		conf, err := di.Get[*config.RuntimeConfig](ctx, config.RuntimeConfigKey)
		if err != nil {
			return nil, err
		}

		s := &service{
			db:     di.MustGet[*gorm.DB](ctx, db.Key),
			logger: di.MustGet[*mylog.Logger](ctx, mylog.Key),
			config: conf,
		}

		if conf.LocalToolAutoMigrate || c.Env == di.EnvTest {
			if err := s.InitializeTools(ctx); err != nil {
				return nil, err
			}
		}

		return s, nil
	})
}
