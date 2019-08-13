package hooks

import (
	"math"

	ping "github.com/sparrc/go-ping"
	"github.com/spf13/cobra"

	"github.com/vivek-26/ipv/reporter"
)

const targetIP = "8.8.8.8" // Google DNS Server IP
const nPing = 5            // Number of times to ping

// PersistentPreRun performs internet connection check
func PersistentPreRun(cmd *cobra.Command, args []string) {
	reporter.Info("Checking internet connection...")
	pinger, err := ping.NewPinger(targetIP)
	if err != nil {
		reporter.Error(err)
	}

	pinger.Count = nPing
	pinger.SetPrivileged(true)
	pinger.Run() // Blocks until finished
	stats := pinger.Statistics()

	// `pinger.Run()` does not return errors. Hence, to know if an error
	// has occurred, we use `pinger.Statistics()`.
	// To determine errors, two conditions are used -
	// 1) `stats.PacketLoss` is NaN
	// 2) Average round trip time (stats.AvgRtt) is 0
	if math.IsNaN(stats.PacketLoss) || stats.AvgRtt == 0 {
		reporter.Error("Internet connection check failed ✗")
	}

	reporter.Success("Internet connection check successful ✓")
}
