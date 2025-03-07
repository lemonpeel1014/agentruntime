package tool_test

import (
	"context"
	"fmt"
	"github.com/habiliai/agentruntime/di"
	"github.com/habiliai/agentruntime/tool"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetWeather(t *testing.T) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		t.Skip("OPENWEATHER_API_KEY 환경 변수가 설정되지 않았습니다")
	}

	ctx := di.WithContainer(context.TODO(), di.EnvTest)

	s := di.MustGet[tool.Manager](ctx, tool.ManagerKey)
	weatherSummary, err := s.GetWeather(ctx, &tool.GetWeatherRequest{
		Location: "Seoul",
		Date:     "2023-10-01",
		Unit:     "c",
	})
	require.NoError(t, err)

	t.Logf("contents: %v", weatherSummary)

	// 3. 출력
	fmt.Printf("🌡️ 최고 기온: %.2f°C\n", weatherSummary.Temperature.Max)
	fmt.Printf("🌡️ 최저 기온: %.2f°C\n", weatherSummary.Temperature.Min)
	fmt.Printf("🌡️ 오후 기온(12:00): %.2f°C\n", weatherSummary.Temperature.Afternoon)
	fmt.Printf("🌡️ 아침 기온(06:00): %.2f°C\n", weatherSummary.Temperature.Morning)
	fmt.Printf("🌡️ 저녁 기온(18:00): %.2f°C\n", weatherSummary.Temperature.Evening)
	fmt.Printf("🌡️ 밤 기온(00:00): %.2f°C\n", weatherSummary.Temperature.Night)
	fmt.Printf("💧 오후 습도: %.2f\n", weatherSummary.Humidity.Afternoon)
	fmt.Printf("🌬️ 최대 풍속: %.2fm/s (방향: %.2f°)\n", weatherSummary.Wind.Max.Speed, weatherSummary.Wind.Max.Direction)
}
