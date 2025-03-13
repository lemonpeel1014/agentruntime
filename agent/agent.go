package agent

import "github.com/habiliai/agentruntime/entity"

func (a *Agent) assignFromEntity(agent *entity.Agent) {
	a.Id = uint32(agent.ID)
	a.Name = agent.Name
	a.ModelName = agent.ModelName
	a.Busy = agent.Busy
	a.Metadata = agent.Metadata.Data()
	a.Role = agent.Role
}
