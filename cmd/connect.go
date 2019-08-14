package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vivek-26/ipv/hooks"
)

// connectCmd represents the connect command
func connectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "List top 5 servers (by latency) and connects to chosen server",
		Long: `  List top 5 servers (by latency) in the chosen country.
		It allows the user to select from the top 5 servers
		and connect to it.
		By default, if no flags are passed, it reads configuration values.
		However, they can be overridden by passing flags. For more details,
		use 'help' command.`,
		PersistentPreRun: hooks.PersistentPreRun,
		PreRun:           hooks.PreRun,
		Run:              func(cmd *cobra.Command, args []string) {},
	}
}
