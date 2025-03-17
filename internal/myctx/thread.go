package myctx

import (
	"context"
	"github.com/habiliai/agentruntime/entity"
)

type contextKeyType string

var (
	contextKey = contextKeyType("ctx.thread")
)

func WithThread(ctx context.Context, thread *entity.Thread) context.Context {
	return context.WithValue(ctx, contextKey, thread)
}

func GetThreadFromContext(ctx context.Context) *entity.Thread {
	thread, ok := ctx.Value(contextKey).(*entity.Thread)
	if !ok {
		return nil
	}
	return thread
}
