package tool_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	di "github.com/habiliai/agentruntime/internal/di"
	"github.com/habiliai/agentruntime/tool"
	"github.com/stretchr/testify/require"
)

func TestPostToX(t *testing.T) {
	consumerKey := os.Getenv("X_CONSUMER_KEY")
	if consumerKey == "" {
		t.Skip("X_CONSUMER_KEY 환경 변수가 설정되지 않았습니다")
	}
	consumerSecret := os.Getenv("X_CONSUMER_SECRET")
	if consumerSecret == "" {
		t.Skip("X_CONSUMER_SECRET 환경 변수가 설정되지 않았습니다")
	}
	accessToken := os.Getenv("X_ACCESS_TOKEN")
	if accessToken == "" {
		t.Skip("X_ACCESS_TOKEN 환경 변수가 설정되지 않았습니다")
	}
	accessTokenSecret := os.Getenv("X_ACCESS_TOKEN_SECRET")
	if accessTokenSecret == "" {
		t.Skip("X_ACCESS_TOKEN_SECRET 환경 변수가 설정되지 않았습니다")
	}
	accountId := os.Getenv("X_ACCOUNT_ID")
	if accountId == "" {
		t.Skip("X_ACCOUNT_ID 환경 변수가 설정되지 않았습니다")
	}

	ctx := di.WithContainer(context.TODO(), di.EnvTest)

	s := di.MustGet[tool.Manager](ctx, tool.ManagerKey)
	content, err := s.PostToX(ctx, &tool.PostToXRequest{
		Text: "Post to X tool test completed.\npost time: " + time.Now().Format(time.RFC3339),
	})
	require.NoError(t, err)

	fmt.Printf("post content: %v\n", content)
}
