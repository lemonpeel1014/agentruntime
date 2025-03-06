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
	"gorm.io/datatypes"
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
	}
)

func (s *service) Run(
	ctx context.Context,
	threadId uint,
	agentId uint,
) error {
	ctx, sess := db.OpenSession(ctx, s.db)

	var thread entity.Thread
	if err := sess.First(&thread, threadId).Error; err != nil {
		return errors.Wrapf(err, "failed to find thread")
	}

	var agent entity.Agent
	if err := sess.First(&agent, agentId).Error; err != nil {
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

	// construct inst values
	instValues := ChatInstValues{
		Agent:               agent,
		RecentConversations: make([]Conversation, 0, len(messages)),
		AvailableActions:    make([]AvailableAction, 0, len(agent.Functions)),
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
	for _, function := range agent.Functions {
		instValues.AvailableActions = append(instValues.AvailableActions, AvailableAction{
			Action:      function.Name,
			Description: function.Description,
		})
	}

	var instBuilder strings.Builder
	if err := chatInstTmpl.Execute(&instBuilder, instValues); err != nil {
		return errors.Wrapf(err, "failed to execute template")
	}

	model := openai.Model(agent.ModelName)

	resp, err := ai.Generate(
		ctx,
		model,
		ai.WithSystemPrompt(agent.System),
		ai.WithTextPrompt(instBuilder.String()),
		ai.WithConfig(&ai.GenerationCommonConfig{
			Temperature: 0.2,
			TopP:        0.5,
			TopK:        16,
		}),
		ai.WithOutputFormat(ai.OutputFormatJSON),
		ai.WithOutputSchema(&Conversation{}),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to generate")
	}

	var conversation Conversation
	if err := json.Unmarshal([]byte(resp.Text()), &conversation); err != nil {
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
}
