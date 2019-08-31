package ipvanish

import (
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
	"github.com/vivek-26/ipv/reporter"
)

const (
	maxRTT = time.Second * 3 // MaxRTT for one ping cycle
	factor = 25              // Used to compute ping cycles
)

// PingAllServers records latencies (RTT) for all ipvanish servers.
// When the number of servers passed is large, ping might not complete
// for all servers. Hence, ping is run multiple times. Whenever, maxRTT
// is passed, OnIdle is called where the server IPs for which latencies
// have been recorded are removed. And ping runs again. factor is used to
// compute the number of ping cycles. It is assumed that in one cycle
// atleast 25 servers are being pinged successfully.
func PingAllServers(servers *[]IPVServer) *[]IPVServer {
	ping := fastping.NewPinger()
	_, err := ping.Network("udp")
	if err != nil {
		reporter.Error(err)
	}

	ping.MaxRTT = maxRTT

	// Add IP addresses to pinger instance
	results := make(map[string]*IPVServer)
	for i := range *servers {
		results[(*servers)[i].IP] = &((*servers)[i])
		_ = ping.AddIP((*servers)[i].IP)
	}

	counter := 0 // Keep track of completed pings

	ping.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		counter++
		if rtt > 0 {
			server := results[addr.String()]
			server.Latency = rtt
		}
	}

	ping.OnIdle = func() {
		// Remove server IPs from pinger for which ping has completed
		for i := range *servers {
			if (*servers)[i].Latency != 0 {
				_ = ping.RemoveIP((*servers)[i].IP)
			}
		}
	}

	pingCycles := len(*servers) / factor
	for i := 0; i <= pingCycles; i++ {
		if counter < len(*servers) {
			err := ping.Run() // Blocking call
			if err != nil {
				reporter.Error(err)
			}
		} else {
			break
		}
	}

	return servers
}
