package agentruntime

import (
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/internal/di"
	"github.com/habiliai/agentruntime/runtime"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <thread_id> <agent_name>",
		Short: "Run agent runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) != 2 {
				return errors.Errorf("thread_id and agent_name are required")
			}

			threadId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.Wrapf(err, "thread_id must be a number")
			}

			agentName := args[1]

			runtime, err := di.Get[runtime.Runtime](ctx, runtime.Key)
			if err != nil {
				return err
			}

			agentManager, err := di.Get[agent.Manager](ctx, agent.ManagerKey)
			if err != nil {
				return err
			}

			ag, err := agentManager.FindAgentByName(ctx, agentName)
			if err != nil {
				return err
			}

			if err := runtime.Run(ctx, uint(threadId), []uint{ag.ID}); err != nil {
				return err
			}

			return nil
		},
	}
}
