package ipvanish

import (
	"fmt"
	"time"
)

// IPVServer represents Ipvanish Server
type IPVServer struct {
	IP       string
	Hostname string
	Latency  time.Duration
}

// ByLatency implements sort.Interface based on Latency field
type ByLatency []IPVServer

// Len returns number of servers
func (s ByLatency) Len() int { return len(s) }

// Less compares latency of two servers
func (s ByLatency) Less(i, j int) bool { return s[i].Latency < s[j].Latency }

// Swap exchanges servers to order them
func (s ByLatency) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Implement stringer interface
func (i IPVServer) String() string {
	return fmt.Sprintf(
		"Host: %v, IP: %v, Latency: %v",
		i.Hostname, i.IP, i.Latency,
	)
}

// ClientInfo represents current connection info of user
type ClientInfo struct {
	IPAddress string `json:"ip_address"`
	Location  struct {
		CountryName string  `json:"country_name"`
		Region      string  `json:"region"`
		City        string  `json:"city"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
	} `json:"location"`
	VPN struct {
		Enabled bool `json:"enabled"`
		Secure  bool `json:"secure"`
	} `json:"vpn"`
}
