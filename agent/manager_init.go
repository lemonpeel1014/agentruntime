package agent

import (
	"context"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/habiliai/agentruntime/internal/di"
	"github.com/habiliai/agentruntime/internal/mylog"
	"github.com/habiliai/agentruntime/tool"
	"gorm.io/gorm"
	"log/slog"
)

var (
	ManagerKey = di.NewKey()
)

func init() {
	di.Register(ManagerKey, func(c context.Context, _ di.Env) (any, error) {
		return &manager{
			logger:      di.MustGet[*slog.Logger](c, mylog.Key),
			db:          di.MustGet[*gorm.DB](c, db.Key),
			toolManager: di.MustGet[tool.Manager](c, tool.ManagerKey),
		}, nil
	})
}
