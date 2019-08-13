package hooks

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/vivek-26/ipv/reporter"
)

// PreRun performs DNS check
func PreRun(cmd *cobra.Command, args []string) {
	reporter.Info("Checking DNS...")
	ipRecords, err := net.LookupIP("google.com")
	if err != nil {
		reporter.Error("DNS check failed ✗")
	}

	// One or more IP addresses should be available
	if len(ipRecords) > 0 {
		reporter.Success("DNS check successful ✓")
		return
	}

	reporter.Error("DNS check failed ✗")
}
