package config_test

import (
	"github.com/habiliai/agentruntime/config"
)

func (s *ConfigTestSuite) TestLoadAgentsFromFiles() {
	testFile := "./testdata/test_agent_1.yaml"

	agentConfigs, err := config.LoadAgentsFromFiles([]string{testFile})
	s.Require().NoError(err)
	s.Require().Len(agentConfigs, 1)

	agentConfig := agentConfigs[0]
	s.T().Logf("AgentConfig: %+v", agentConfig)

	s.Require().Equal("Alice", agentConfig.Name)
	s.Require().Equal("USER", agentConfig.MessageExamples[0].Messages[0].Name)
	s.Require().Equal("done_agent", agentConfig.MessageExamples[0].Messages[0].Action)
}
