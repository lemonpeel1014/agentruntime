package agent_test

import "github.com/habiliai/agentruntime/config"

func (s *AgentTestSuite) TestSaveAgent() {
	agentConfigs, err := config.LoadAgentsFromFiles([]string{"testdata/habiliai.agent.yaml"})
	s.Require().NoError(err)

	s.Require().Len(agentConfigs, 1)
	agentConfig := agentConfigs[0]

	s.Require().Len(agentConfig.MessageExamples, 1)

	agent, err := s.manager.SaveAgentFromConfig(s, agentConfig)
	s.Require().NoError(err)

	s.Require().Len(agent.MessageExamples, 1)
}
