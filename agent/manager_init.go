package agent

import (
	"context"
	"github.com/habiliai/agentruntime/di"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/habiliai/agentruntime/internal/mylog"
	"gorm.io/gorm"
	"log/slog"
)

var (
	ManagerKey = di.NewKey()
)

func init() {
	di.Register(ManagerKey, func(c context.Context, _ *di.Container) (any, error) {
		return &manager{
			logger: di.MustGet[*slog.Logger](c, mylog.Key),
			db:     di.MustGet[*gorm.DB](c, db.Key),
		}, nil
	})
}
