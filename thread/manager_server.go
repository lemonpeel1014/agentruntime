package thread

import (
	"context"
	"github.com/habiliai/agentruntime/internal/di"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type managerServer struct {
	UnsafeThreadManagerServer

	manager Manager
}

func (m *managerServer) CreateThread(ctx context.Context, req *CreateThreadRequest) (*CreateThreadResponse, error) {
	thr, err := m.manager.CreateThread(ctx, req.Instruction)
	if err != nil {
		return nil, err
	}

	return &CreateThreadResponse{
		ThreadId: uint32(thr.ID),
	}, nil
}

func (m *managerServer) GetThread(ctx context.Context, request *GetThreadRequest) (*GetThreadResponse, error) {
	thr, err := m.manager.GetThreadById(ctx, uint(request.ThreadId))
	if err != nil {
		return nil, err
	}

	return &GetThreadResponse{
		Thread: &Thread{
			Id:          uint32(thr.ID),
			Instruction: thr.Instruction,
			CreatedAt:   timestamppb.New(thr.CreatedAt),
			UpdatedAt:   timestamppb.New(thr.UpdatedAt),
		},
	}, nil
}

func (m *managerServer) AddMessage(ctx context.Context, request *AddMessageRequest) (*AddMessageResponse, error) {
	msg, err := m.manager.AddMessage(ctx, uint(request.ThreadId), request.Message)
	if err != nil {
		return nil, err
	}

	return &AddMessageResponse{
		MessageId: uint32(msg.ID),
	}, nil
}

var (
	_                ThreadManagerServer = (*managerServer)(nil)
	ManagerServerKey                     = di.NewKey()
)

func init() {
	di.Register(ManagerServerKey, func(ctx context.Context, _ di.Env) (any, error) {
		return &managerServer{
			manager: di.MustGet[Manager](ctx, ManagerKey),
		}, nil
	})
}
