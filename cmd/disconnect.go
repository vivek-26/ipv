package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vivek-26/ipv/ipvanish"
	"github.com/vivek-26/ipv/utils"
)

// disconnectCmd represents disconnect command
func disconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnects from VPN by killing the OpenVPN process",
		Long:  utils.HelpText("Disconnects from VPN by killing the OpenVPN process"),
		Run: func(cmd *cobra.Command, args []string) {
			// Disconnect from vpn
			ipvanish.Disconnect()
		},
	}
}
