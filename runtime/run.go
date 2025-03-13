package runtime

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/firebase/genkit/go/ai"
	"github.com/habiliai/agentruntime/entity"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/pkg/errors"
	"github.com/yukinagae/genkit-go-plugins/plugins/openai"
	"golang.org/x/sync/errgroup"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
	"slices"
	"strings"
	"text/template"
)

var (
	//go:embed data/instructions/chat.md.tmpl
	chatInst     string
	chatInstTmpl = template.Must(template.New("chatInst").Funcs(funcMap()).Parse(chatInst))
)

type (
	Conversation struct {
		User   string `json:"user"`
		Text   string `json:"text"`
		Action string `json:"action"`
	}

	AvailableAction struct {
		Action      string `json:"action"`
		Description string `json:"description"`
	}

	ChatInstValues struct {
		Agent entity.Agent

		RecentConversations []Conversation
		Knowledge           []string
		AvailableActions    []AvailableAction
		MessageExamples     [][]entity.MessageExample
	}
)

func (s *service) Run(
	ctx context.Context,
	threadId uint,
	agentIds []uint,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, sess := db.OpenSession(ctx, s.db)

	var thread entity.Thread
	if err := sess.First(&thread, threadId).Error; err != nil {
		return errors.Wrapf(err, "failed to find thread")
	}

	var agents []entity.Agent
	if err := sess.Preload(clause.Associations).Find(&agents, "id in ?", agentIds).Error; err != nil {
		return errors.Wrapf(err, "failed to find agent")
	}

	var messages []entity.Message
	if err := sess.
		Order("created_at DESC").
		Limit(200).
		Find(&messages, "thread_id = ?", threadId).Error; err != nil {
		return errors.Wrapf(err, "failed to find messages")
	}

	slices.SortStableFunc(messages, func(a, b entity.Message) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		} else if a.CreatedAt.After(b.CreatedAt) {
			return 1
		} else {
			return 0
		}
	})

	var eg errgroup.Group
	for _, agent := range agents {
		eg.Go(func() error {
			// construct inst values
			instValues := ChatInstValues{
				Agent:               agent,
				MessageExamples:     agent.MessageExamples,
				RecentConversations: make([]Conversation, 0, len(messages)),
				AvailableActions:    make([]AvailableAction, 0, len(agent.Tools)),
			}

			// build recent conversations
			for _, msg := range messages {
				instValues.RecentConversations = append(instValues.RecentConversations, Conversation{
					User:   msg.User,
					Text:   msg.Content.Data().Text,
					Action: msg.Content.Data().Action,
				})
			}

			// build available actions
			tools := make([]ai.Tool, 0, len(agent.Tools))
			for _, tool := range agent.Tools {
				instValues.AvailableActions = append(instValues.AvailableActions, AvailableAction{
					Action:      tool.Name,
					Description: tool.Description,
				})
				tools = append(tools, s.toolManager.GetLocalTool(ctx, tool.LocalToolName))
			}

			var promptBuilder strings.Builder
			if err := chatInstTmpl.Execute(&promptBuilder, instValues); err != nil {
				return errors.Wrapf(err, "failed to execute template")
			}
			prompt := promptBuilder.String()

			s.logger.Debug("call agent runtime's run", "prompt", prompt)

			model := openai.Model(agent.ModelName)

			resp, err := ai.Generate(
				ctx,
				model,
				ai.WithCandidates(1),
				ai.WithSystemPrompt(agent.System),
				ai.WithTextPrompt(prompt),
				ai.WithConfig(&ai.GenerationCommonConfig{
					Temperature: 0.2,
					TopP:        0.5,
					TopK:        16,
				}),
				// TODO: Cannot support using tools with output format
				ai.WithOutputFormat(ai.OutputFormatJSON),
				ai.WithOutputSchema(&Conversation{}),
				ai.WithTools(tools...),
			)
			if err != nil {
				return errors.Wrapf(err, "failed to generate")
			}

			responseText := resp.Text()

			var conversation Conversation
			if err := json.Unmarshal([]byte(responseText), &conversation); err != nil {
				return errors.Wrapf(err, "failed to unmarshal conversation")
			}

			newMessage := entity.Message{
				ThreadID: threadId,
				User:     agent.Name,
				Content: datatypes.NewJSONType(entity.MessageContent{
					Text: conversation.Text,
				}),
			}
			if err := sess.Create(&newMessage).Error; err != nil {
				return errors.Wrapf(err, "failed to create message")
			}

			return nil
		})
	}

	return eg.Wait()
}
