package runtime

import (
	"context"
	"github.com/habiliai/agentruntime/di"
	"github.com/mokiat/gog"
)

type agentRuntimeServer struct {
	UnsafeAgentRuntimeServer

	runtime Runtime
}

func (a *agentRuntimeServer) Run(ctx context.Context, req *RunRequest) (*RunResponse, error) {
	err := a.runtime.Run(ctx, uint(req.ThreadId), gog.Map(req.AgentIds, func(id uint32) uint {
		return uint(id)
	}))
	return &RunResponse{}, err
}

var (
	_         AgentRuntimeServer = (*agentRuntimeServer)(nil)
	ServerKey                    = di.NewKey()
)

func init() {
	di.Register(ServerKey, func(c context.Context, _ *di.Container) (any, error) {
		return &agentRuntimeServer{
			runtime: di.MustGet[Runtime](c, Key),
		}, nil
	})
}
