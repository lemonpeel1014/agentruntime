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
		Use:   "run <thread_id> <agent_name> [...<agent_name>]",
		Short: "Run agent runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(args) < 2 {
				return errors.Errorf("thread_id and agent_name are required")
			}

			threadId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.Wrapf(err, "thread_id must be a number")
			}

			agentNames := args[1:]

			runtime, err := di.Get[runtime.Runtime](ctx, runtime.Key)
			if err != nil {
				return err
			}

			agentManager, err := di.Get[agent.Manager](ctx, agent.ManagerKey)
			if err != nil {
				return err
			}

			agentIds := make([]uint, 0, len(agentNames))
			for _, name := range agentNames {
				ag, err := agentManager.FindAgentByName(ctx, name)
				if err != nil {
					return err
				}
				agentIds = append(agentIds, ag.ID)
			}

			if err := runtime.Run(ctx, uint(threadId), agentIds); err != nil {
				return err
			}

			return nil
		},
	}
}
