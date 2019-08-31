package hooks

import (
	"bytes"
	"os/exec"
	"time"

	"github.com/spf13/cobra"

	"github.com/vivek-26/ipv/reporter"
)

const targetAddr = "ipvanish.com"
const maxRTT = time.Second * 1 // Max round trip time

// PersistentPreRun checks for openvpn binary
func PersistentPreRun(cmd *cobra.Command, args []string) {
	checkOpenvpnBinary()
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
