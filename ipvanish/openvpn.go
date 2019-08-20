package ipvanish

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vivek-26/ipv/reporter"
)

const (
	ipvanishVpnConfigURL = "http://files.ipvanish.com/ipvanish-openvpn-config.txt"
	proto                = "PROTOCOL"
	server               = "SERVER"
	credentialsFile      = "credentials.txt"
	pid                  = "PIDFILE"
	pidFile              = "/ipv.pid"
	openvpnFile          = "/openvpn.conf"
)

// GenerateVPNConfig generate OpenVPN config file for hostname
func GenerateVPNConfig(hostname, protocol string) {
	resp, err := http.Get(ipvanishVpnConfigURL)
	if err != nil {
		reporter.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		reporter.Error(err)
	}

	// Generate VPN config
	vpnConfig := string(body)

	// Replace protocol
	vpnConfig = strings.ReplaceAll(vpnConfig, proto, protocol)

	// Replace server hostname
	vpnConfig = strings.ReplaceAll(vpnConfig, server, hostname)

	// Replace credentials file path
	vpnConfig = strings.ReplaceAll(vpnConfig, credentialsFile, getCredentialsFilepath())

	// Replace PID file path
	vpnConfig = strings.ReplaceAll(vpnConfig, pid, getPIDFilepath())

	// Write VPN config
	vpnConfigFile := getVPNConfigFilepath()
	if err := ioutil.WriteFile(vpnConfigFile, []byte(vpnConfig), 0644); err != nil {
		reporter.Error(err)
	}
}

// getPIDFilepath returns openvpn pid file path
func getPIDFilepath() string {
	return getConfigDirPath() + pidFile
}

// getVPNConfigFilepath returns openvpn pid file path
func getVPNConfigFilepath() string {
	return getConfigDirPath() + openvpnFile
}
