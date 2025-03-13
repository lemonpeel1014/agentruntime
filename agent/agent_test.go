package agent_test

import (
	"context"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/internal/di"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AgentTestSuite struct {
	suite.Suite
	context.Context

	manager agent.Manager
}

func (s *AgentTestSuite) SetupTest() {
	s.Context = context.TODO()
	s.Context = di.WithContainer(s.Context, di.EnvTest)

	s.manager = di.MustGet[agent.Manager](s.Context, agent.ManagerKey)
}

func TestAgents(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}
