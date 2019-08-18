package hooks

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/tatsushid/go-fastping"

	"github.com/vivek-26/ipv/reporter"
)

const targetAddr = "ipvanish.com"
const maxRTT = time.Second * 2 // Max round trip time

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

	onRecv, onIdle := make(chan struct{}), make(chan struct{})

	// Received ICMP message handler
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		onRecv <- struct{}{}
	}

	// Max RTT expiration handler
	p.OnIdle = func() {
		onIdle <- struct{}{}
	}

	p.RunLoop() // Non blocking

	c := make(chan os.Signal, 1) // Look for `ctrl-c`
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	select {
	case <-c:
		reporter.Error("Process stopped by user")
	case <-onRecv:
		reporter.Success("Internet connection check successful ✓")
	case <-onIdle:
		reporter.Error("Internet connection check failed ✗")
	case <-p.Done():
		if err := p.Err(); err != nil {
			reporter.Error(err)
		}
	}

	signal.Stop(c)
	p.Stop()
}
