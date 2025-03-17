package tool

import (
	"context"
	"github.com/firebase/genkit/go/ai"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/entity"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/habiliai/agentruntime/internal/mylog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	service struct {
		logger *mylog.Logger
		db     *gorm.DB
		config *config.RuntimeConfig
	}
)

func (s *service) InitializeTools(ctx context.Context) error {
	_, tx := db.OpenSession(ctx, s.db)
	return tx.Transaction(func(tx *gorm.DB) error {
		for _, toolName := range localToolNames {
			tool := ai.LookupTool(toolName).Definition()
			if err := tx.Save(&entity.Tool{
				Name:          toolName,
				Description:   tool.Description,
				LocalToolName: tool.Name,
			}).Error; err != nil {
				return errors.Wrapf(err, "failed to save tool")
			}
		}

		return nil
	})
}

func (s *service) GetLocalTool(_ context.Context, toolName string) ai.Tool {
	return ai.LookupTool(toolName)
}

func (s *service) GetTools(ctx context.Context, names []string) ([]entity.Tool, error) {
	_, tx := db.OpenSession(ctx, s.db)

	var tools []entity.Tool
	if err := tx.Find(&tools, "name IN ?", names).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find tools")
	}

	return tools, nil
}

var (
	_              Manager = (*service)(nil)
	localToolNames []string
)

func RegisterLocalTool[In any, Out any](name string, description string, fn func(context.Context, In) (Out, error)) ai.Tool {
	localToolNames = append(localToolNames, name)
	return ai.DefineTool(
		name,
		description,
		func(ctx context.Context, input In) (Out, error) {
			out, err := fn(ctx, input)
			if err == nil {
				appendCallData(ctx, CallData{
					Name:      name,
					Arguments: input,
					Result:    out,
				})
			}
			return out, err
		},
	)
}
