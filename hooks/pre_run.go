package hooks

import (
	"net"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/tatsushid/go-fastping"

	"github.com/vivek-26/ipv/reporter"
)

// PreRun performs internet and dns check
func PreRun(cmd *cobra.Command, args []string) {
	internetCheck()
	dnsCheck()
}

// internetCheck verifies internet connectivity
func internetCheck() {
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

// dnsCheck verifies the dns service
func dnsCheck() {
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
