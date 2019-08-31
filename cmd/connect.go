package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/briandowns/spinner"
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

			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			_ = s.Color("yellow", "bold")
			r := &reporter.Spinner{Spin: s}
			r.Info("Updating servers list")

			// Get all ipvanish servers for given country
			countryCode := viper.GetString("countryCode")
			servers := ipvanish.GetServers(countryCode)
			if len(*servers) > 0 {
				r.Success()
			} else {
				r.Error()
			}

			// Ping all servers
			r.Info(
				fmt.Sprintf("Pinging %v servers in %v", len(*servers), countryCode),
			)
			servers = ipvanish.PingAllServers(servers)

			// Sort based on latency
			sort.Sort(ipvanish.ByLatency(*servers))
			r.Success()

			// Ask user to choose server
			hostname := ipvanish.SelectServerPrompt(servers, 5)

			// Create credentials file
			username := viper.GetString("username")
			password := viper.GetString("password")
			ipvanish.CreateCredentials(username, password)

			// Generate VPN config
			protocol := viper.GetString("protocol")
			ipvanish.GenerateVPNConfig(hostname, protocol)

			// Connect to vpn server
			ipvanish.Connect()
			reporter.Info("Connected ✔")
		},
	}
}
