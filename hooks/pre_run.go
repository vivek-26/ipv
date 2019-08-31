package hooks

import (
	"net"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/vivek-26/ipv/reporter"
)

// PreRun performs DNS check
func PreRun(cmd *cobra.Command, args []string) {
	// Create and start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	_ = s.Color("yellow", "bold")
	r := &reporter.Spinner{Spin: s}
	r.Info("Checking DNS")

	ipRecords, err := net.LookupIP("google.com")
	if err != nil {
		reporter.Error("DNS check failed ✗")
	}

	// One or more IP addresses should be available
	if len(ipRecords) > 0 {
		r.Success()
		return
	}

	r.Error()
}
