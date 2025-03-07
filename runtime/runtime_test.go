package runtime_test

import (
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/di"
	"github.com/habiliai/agentruntime/internal/mytesting"
	"github.com/habiliai/agentruntime/runtime"
	"github.com/habiliai/agentruntime/thread"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type AgentRuntimeTestSuite struct {
	mytesting.Suite

	agents        []config.AgentConfig
	runtime       runtime.Runtime
	agentManager  agent.Manager
	threadManager thread.Manager
}

func (s *AgentRuntimeTestSuite) SetupTest() {
	os.Setenv("ENV_TEST_FILE", "../.env.test")
	s.Suite.SetupTest()

	var err error

	s.agents, err = config.LoadAgentsFromFiles([]string{"./testdata/test1.agent.yaml"})
	s.Require().NoError(err)

	s.runtime = di.MustGet[runtime.Runtime](s, runtime.Key)
	s.agentManager = di.MustGet[agent.Manager](s, agent.ManagerKey)
	s.threadManager = di.MustGet[thread.Manager](s, thread.ManagerKey)

	s.Require().NoError(err)
}

func (s *AgentRuntimeTestSuite) TearDownTest() {
	defer s.Suite.TearDownTest()
}

func TestAgentRuntime(t *testing.T) {
	suite.Run(t, new(AgentRuntimeTestSuite))
}
