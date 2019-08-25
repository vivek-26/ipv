package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vivek-26/ipv/ipvanish"
)

// disconnectCmd represents disconnect command
func disconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnects from VPN by killing VPN process",
		Long:  "  Disconnects from VPN by killing VPN process",
		Run: func(cmd *cobra.Command, args []string) {
			// Disconnect from vpn
			ipvanish.Disconnect()
		},
	}
}