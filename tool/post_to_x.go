package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/habiliai/agentruntime/internal/di"
	"github.com/pkg/errors"
)

type (
	PostToXRequest struct {
		Text string `json:"text" jsonschema:"required,description=Text to post"`
	}

	PostToXResponse struct {
		Link       string `json:"link"`
		Content    string `json:"content"`
		UploadedAt string `json:"uploaded_at"`
	}
)

func callPostApi(ctx context.Context, consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string, accountId string, message string) (*PostToXResponse, error) {
	type (
		PostApiResponse struct {
			Data struct {
				Id   string `json:"id"`
				Text string `json:"text"`
			} `json:"data"`
			Errors []struct {
				Detail string `json:"detail"`
				Title  string `json:"title"`
				Type   string `json:"type"`
				Status int    `json:"status"`
			} `json:"errors"`
		}

		PostApiPayload struct {
			Text string `json:"text"`
		}
	)
	const tweetUrl = "https://api.x.com/2/tweets"

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(ctx, token)

	payload, err := json.Marshal(&PostApiPayload{
		Text: message,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload")
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", tweetUrl, bytes.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to post tweet")
	}
	defer res.Body.Close()

	responseBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read response body")
	}
	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("failed to post tweet, status code: %d, body: '%s'", res.StatusCode, string(responseBodyBytes))
	}

	var resp PostApiResponse
	if err := json.Unmarshal(responseBodyBytes, &resp); err != nil {
		return nil, errors.Wrapf(err, "failed to decode response body")
	}

	if len(resp.Errors) > 0 {
		return nil, errors.Errorf("failed to post tweet, errors: %v", resp.Errors)
	}

	postId := resp.Data.Id

	return &PostToXResponse{
		Link:       "https://x.com/" + accountId + "/status/" + postId,
		Content:    resp.Data.Text,
		UploadedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *service) PostToX(ctx context.Context, req *PostToXRequest) (*PostToXResponse, error) {
	if req.Text == "" {
		return nil, errors.New("Text is required")
	}

	postResult, err := callPostApi(
		ctx,
		s.config.XConsumerKey,
		s.config.XConsumerSecret,
		s.config.XAccessToken,
		s.config.XAccessTokenSecret,
		s.config.XAccountId,
		req.Text,
	)

	return postResult, err
}

func init() {
	RegisterLocalTool(
		"post_to_x",
		"Post to X when you need it",
		func(ctx context.Context, req struct {
			*PostToXRequest
		}) (res struct {
			*PostToXResponse
		}, err error) {
			s := di.MustGet[*service](ctx, ManagerKey)
			res.PostToXResponse, err = s.PostToX(ctx, req.PostToXRequest)
			return
		},
	)
}
