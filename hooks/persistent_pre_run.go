package hooks

import (
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/tatsushid/go-fastping"

	"github.com/vivek-26/ipv/reporter"
)

const targetAddr = "ipvanish.com"
const maxRTT = time.Second * 1 // Max round trip time

// PersistentPreRun performs internet connection check
func PersistentPreRun(cmd *cobra.Command, args []string) {
	reporter.Info("Checking internet connection...")
	p := fastping.NewPinger()
	_, err := p.Network("udp")
	if err != nil {
		reporter.Error(err)
	}

	p.MaxRTT = maxRTT

	ra, err := net.ResolveIPAddr("ip4:icmp", targetAddr)
	if err != nil {
		reporter.Error(err)
	}

	// Add target IP address
	p.AddIPAddr(ra)

	var isHostReachable bool

	// Received ICMP message handler
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		if rtt > 0 {
			isHostReachable = true
		}
	}

	// Max RTT expiration handler
	p.OnIdle = func() {
		if isHostReachable {
			reporter.Success("Internet connection check successful ✓")
		} else {
			reporter.Error("Internet connection check failed ✗")
		}
	}

	err = p.Run() // Blocking
	if err != nil {
		reporter.Error("Internet connection check failed ✗")
	}
}
