package thread

import (
	"context"
	"github.com/habiliai/agentruntime/entity"
	"github.com/habiliai/agentruntime/internal/di"
	"github.com/mokiat/gog"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type managerServer struct {
	UnsafeThreadManagerServer

	manager Manager
}

func (m *managerServer) GetMessages(req *GetMessagesRequest, stream ThreadManager_GetMessagesServer) error {
	ctx := stream.Context()
	cursor := uint(0)
	order := "ASC"
	if req.Order == GetMessagesRequest_LATEST {
		order = "DESC"
	}
	for {
		messages, err := m.manager.GetMessages(ctx, uint(req.ThreadId), order, cursor, uint(req.Limit))
		if err != nil {
			return err
		}
		if len(messages) == 0 {
			break
		}

		resp := &GetMessagesResponse{
			Messages: gog.Map(messages, func(msg entity.Message) *Message {
				return &Message{
					Id:        uint32(msg.ID),
					Content:   msg.Content.Data().Text,
					CreatedAt: timestamppb.New(msg.CreatedAt),
					UpdatedAt: timestamppb.New(msg.UpdatedAt),
					Sender:    msg.User,
				}
			}),
		}
		if err := stream.Send(resp); err != nil {
			return errors.Wrapf(err, "failed to send messages")
		}
		cursor = messages[len(messages)-1].ID
	}

	return nil
}

func (m *managerServer) GetNumMessages(ctx context.Context, req *GetNumMessagesRequest) (*GetNumMessagesResponse, error) {
	numMessages, err := m.manager.GetNumMessages(ctx, uint(req.ThreadId))
	if err != nil {
		return nil, err
	}

	return &GetNumMessagesResponse{
		NumMessages: uint32(numMessages),
	}, nil
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
