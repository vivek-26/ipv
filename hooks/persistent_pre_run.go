package hooks

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/vivek-26/ipv/ipvanish"
	"github.com/vivek-26/ipv/reporter"
)

const targetAddr = "ipvanish.com"
const maxRTT = time.Second * 1 // Max round trip time

// PersistentPreRun checks for openvpn binary
func PersistentPreRun(cmd *cobra.Command, args []string) {
	checkOpenvpnBinary()

	// Check if a vpn process is already running, if true, then abort
	isRunning, pid, _ := ipvanish.IsVPNProcessRunning()
	if isRunning {
		reporter.Error(
			strings.Join(
				[]string{
					fmt.Sprintf("A VPN process with PID %v is already running.", pid),
					"Please disconnect using command `sudo ipv disconnect`",
				},
				"\n"),
		)
	}
}

// checkOpenvpnBinary makes sure that openvpn binary is installed on the system
func checkOpenvpnBinary() {
	cmd := exec.Command(
		"command", "-v", "openvpn",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		reporter.Error(
			"OpenVPN package is not installed or openvpn binary is not in standard PATH",
		)
	}

	if stderr.String() != "" {
		reporter.Error(stderr.String())
	}
}
