package myctx

import (
	"context"
	"github.com/habiliai/agentruntime/entity"
)

type agentKeyType string

var agentKey = agentKeyType("ctx.agent")

func WithAgent(ctx context.Context, agent *entity.Agent) context.Context {
	return context.WithValue(ctx, agentKey, agent)
}

func GetAgentFromContext(ctx context.Context) *entity.Agent {
	agent, ok := ctx.Value(agentKey).(*entity.Agent)
	if !ok {
		return nil
	}
	return agent
}
