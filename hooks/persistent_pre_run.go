package hooks

import (
	"net"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/tatsushid/go-fastping"

	"github.com/vivek-26/ipv/reporter"
)

const targetAddr = "ipvanish.com"
const maxRTT = time.Second * 1 // Max round trip time

// PersistentPreRun performs internet connection check
func PersistentPreRun(cmd *cobra.Command, args []string) {
	// Create and start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	_ = s.Color("yellow", "bold")
	r := &reporter.Spinner{Spin: s}
	r.Info("Checking internet connection")

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
			r.Success()
		} else {
			r.Error()
		}
	}

	err = p.Run() // Blocking
	if err != nil {
		reporter.Error("Internet connection check failed ✗")
	}
}
