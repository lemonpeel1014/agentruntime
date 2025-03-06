package agentruntime

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agentruntime",
		Short: "Agent runtime by HabiliAI",
	}

	cmd.AddCommand(
		newRunCmd(),
		newAgentCmd(),
		newThreadCmd(),
	)

	return cmd
}
