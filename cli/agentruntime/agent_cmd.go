package agentruntime

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/di"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agent",
		Short:   "Agent commands",
		Aliases: []string{"agents"},
	}

	listAgentsCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List agents",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()

				agentManager, err := di.Get[agent.Manager](ctx, agent.ManagerKey)
				if err != nil {
					return err
				}

				var (
					cursor uint = 0
					limit  uint = 10
				)

				screen, err := tcell.NewScreen()
				if err != nil {
					return err
				}
				if err := screen.Init(); err != nil {
					return err
				}
				defer screen.Fini()

				return listScreen(ctx, screen, ListScreenRequest{}, func() (messages []string, err error) {
					agents, err := agentManager.GetAgents(ctx, cursor, limit)
					if err != nil {
						return nil, err
					}
					if len(agents) == 0 {
						return
					}

					for _, agent := range agents {
						msg := fmt.Sprintf("Agent: %s\n", agent.Name)
						messages = append(messages, msg)
					}
					cursor = agents[len(agents)-1].ID

					return
				})
			},
		}

		return cmd
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "create <agent-config-file> [...<agent-config-file>]",
			Short: "Create agents from config files",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				if len(args) < 1 {
					return errors.Errorf("agent-config-file is required at least once")
				}

				agentManager, err := di.Get[agent.Manager](ctx, agent.ManagerKey)
				if err != nil {
					return err
				}

				agentConfigs, err := config.LoadAgentsFromFiles(args)
				if err != nil {
					return err
				}

				for _, agentConfig := range agentConfigs {
					if _, err := agentManager.SaveAgentFromConfig(ctx, agentConfig); err != nil {
						return err
					}
				}

				return nil
			},
		},
		listAgentsCmd(),
	)

	return cmd
}
