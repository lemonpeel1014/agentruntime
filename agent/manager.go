package agent

import (
	"context"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/entity"
	myerrors "github.com/habiliai/agentruntime/errors"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/habiliai/agentruntime/tool"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"log/slog"
)

type (
	Manager interface {
		FindAgentByName(ctx context.Context, name string) (*entity.Agent, error)
		SaveAgentFromConfig(ctx context.Context, ac config.AgentConfig) (entity.Agent, error)
		GetAgents(ctx context.Context, cursor uint, limit uint) ([]entity.Agent, error)
		GetAgent(ctx context.Context, id uint) (*entity.Agent, error)
	}
	manager struct {
		logger      *slog.Logger
		db          *gorm.DB
		toolManager tool.Manager
	}
)

var (
	_ Manager = (*manager)(nil)
)

func (s *manager) GetAgent(ctx context.Context, id uint) (*entity.Agent, error) {
	_, tx := db.OpenSession(ctx, s.db)

	var agent entity.Agent
	if err := tx.First(&agent, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find agent")
	}

	return &agent, nil
}

func (s *manager) GetAgents(ctx context.Context, cursor uint, limit uint) ([]entity.Agent, error) {
	_, tx := db.OpenSession(ctx, s.db)

	var agents []entity.Agent
	if err := tx.Where("id > ?", cursor).Order("id ASC").Limit(int(limit)).Find(&agents).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find agents")
	}

	return agents, nil
}

func (s *manager) FindAgentByName(ctx context.Context, name string) (*entity.Agent, error) {
	_, tx := db.OpenSession(ctx, s.db)

	var agent entity.Agent
	if r := tx.Find(&agent, "name ILIKE ?", name); r.Error != nil {
		return nil, errors.Wrapf(r.Error, "failed to find agent")
	} else if r.RowsAffected == 0 {
		return nil, errors.Wrapf(myerrors.ErrNotFound, "failed to find agent")
	}

	return &agent, nil
}

func (s *manager) SaveAgentFromConfig(
	ctx context.Context,
	ac config.AgentConfig,
) (agent entity.Agent, err error) {
	_, tx := db.OpenSession(ctx, s.db)
	if err := tx.Find(&agent, "name = ?", ac.Name).Error; err != nil {
		return agent, errors.Wrapf(err, "failed to find agent")
	}

	agent.Name = ac.Name
	agent.System = ac.System
	agent.Bio = ac.Bio
	agent.Role = ac.Role
	agent.Lore = ac.Lore
	agent.MessageExamples = make([][]entity.MessageExample, 0, len(ac.MessageExamples))
	agent.ModelName = ac.Model
	if agent.ModelName == "" {
		agent.ModelName = "gpt-4o"
	}
	for _, ex := range ac.MessageExamples {
		var messages []entity.MessageExample
		for _, msg := range ex.Messages {
			messages = append(messages, entity.MessageExample{
				User:    msg.Name,
				Text:    msg.Text,
				Actions: msg.Actions,
			})
		}
		agent.MessageExamples = append(agent.MessageExamples, messages)
	}
	tools, err := s.toolManager.GetTools(ctx, ac.Tools)
	if err != nil {
		return agent, err
	}
	agent.Tools = tools
	agent.Metadata = datatypes.NewJSONType(ac.Metadata)
	agent.Knowledge = ac.Knowledge

	if err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&agent).Error; err != nil {
			return errors.Wrapf(err, "failed to save agent")
		}
		return nil
	}); err != nil {
		return agent, err
	}

	return agent, nil
}
