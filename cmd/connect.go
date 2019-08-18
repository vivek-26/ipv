package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vivek-26/ipv/hooks"
	"github.com/vivek-26/ipv/ipvanish"
	"github.com/vivek-26/ipv/reporter"
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
		Run: func(cmd *cobra.Command, args []string) {
			// Print current IP address
			connInfo := ipvanish.GetClientInfo()
			reporter.Success(
				fmt.Sprintf("Your current IP address: %v", connInfo.IPAddress),
			)

			// Get all ipvanish servers for given country
			reporter.Info("Updating server list...")
			countryCode := viper.GetString("countryCode")
			servers := ipvanish.GetServers(countryCode)

			// Ping all servers
			reporter.Info("Pinging all servers...")
			servers = ipvanish.PingAllServers(servers)

			// Sort based on latency
			sort.Sort(ipvanish.ByLatency(*servers))

			// Ask user to choose server
			_ = ipvanish.SelectServerPrompt(servers, 5)
		},
	}
}
