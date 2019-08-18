package ipvanish

import (
	"time"
)

// IPVServer represents Ipvanish Server
type IPVServer struct {
	IP       string
	Hostname string
	Latency  time.Duration
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
