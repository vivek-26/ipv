package ipvanish

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/vivek-26/ipv/reporter"
)

// Connect connects to selected ipvanish vpn server
func Connect() {
	vpnCfgFilePath := getVPNConfigFilepath()
	cmd := exec.Command(
		"sudo", "openvpn", "--config", vpnCfgFilePath,
	)

	reporter.Info(
		fmt.Sprintf("Executing: %s", strings.Join(cmd.Args, " ")),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		reporter.Error(err)
	}

	reporter.Info(string(out))
}
