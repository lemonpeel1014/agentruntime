package agent

import (
	"context"
	"github.com/habiliai/agentruntime/internal/di"
)

type managerServer struct {
	UnsafeAgentManagerServer

	manager Manager
}

func (m *managerServer) GetAgent(ctx context.Context, req *GetAgentRequest) (*Agent, error) {
	agent, err := m.manager.GetAgent(ctx, uint(req.AgentId))
	if err != nil {
		return nil, err
	}

	return &Agent{
		Id:        uint32(agent.ID),
		Name:      agent.Name,
		ModelName: agent.ModelName,
		Busy:      agent.Busy,
		Metadata:  agent.Metadata.Data(),
	}, nil
}

func (m *managerServer) GetAgentByName(ctx context.Context, request *GetAgentByNameRequest) (*Agent, error) {
	agent, err := m.manager.FindAgentByName(ctx, request.Name)
	if err != nil {
		return nil, err
	}

	return &Agent{
		Id:        uint32(agent.ID),
		Name:      agent.Name,
		ModelName: agent.ModelName,
		Busy:      agent.Busy,
		Metadata:  agent.Metadata.Data(),
	}, nil
}

var (
	_                AgentManagerServer = (*managerServer)(nil)
	ManagerServerKey                    = di.NewKey()
)

func init() {
	di.Register(ManagerServerKey, func(c context.Context, _ di.Env) (any, error) {
		return &managerServer{
			manager: di.MustGet[Manager](c, ManagerKey),
		}, nil
	})
}
