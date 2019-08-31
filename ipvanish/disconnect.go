package ipvanish

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	ps "github.com/mitchellh/go-ps"
	"github.com/vivek-26/ipv/reporter"
)

const openvpnProcessName = "openvpn"

// Disconnect kills the vpn process
func Disconnect() {
	reporter.Info("Disconnecting from VPN...")

	isRunning, pid, err := IsVPNProcessRunning()
	if err != nil {
		reporter.Error(err)
	}

	if !isRunning {
		reporter.Error(
			fmt.Sprintf("Process with ID %v is not running", pid),
		)
	}

	process, err := ps.FindProcess(pid)
	if err != nil {
		reporter.Error(err)
	}

	processName := process.Executable()
	if processName != openvpnProcessName {
		reporter.Error("Invalid PID")
	}

	killProcess(pid) // kill vpn process
}

// IsVPNProcessRunning checks if VPN process is running.
// It returns a boolean, an int (pid) and error (if any).
func IsVPNProcessRunning() (bool, int, error) {
	// Find Process ID
	pidFilePath := getPIDFilepath()
	pidF, err := os.Open(pidFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, 0, errors.New("cannot find pid file")
		}
		return false, 0, err
	}

	pidFData, err := ioutil.ReadAll(pidF)
	if err != nil {
		return false, 0, errors.New("cannot read pid file data")
	}

	pidS := strings.TrimSpace(string(pidFData))
	pid, err := strconv.Atoi(pidS) // process id
	if err != nil {
		return false, 0, errors.New("bad process id")
	}

	// Check if VPN Process is running
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false, 0, err
	}

	// Double check if process is running and alive
	// by sending a signal 0
	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		if err == syscall.ESRCH {
			return false, 0, nil
		}
		return false, 0, err
	}

	return true, pid, nil
}

// killProcess kills a process using its `pid`
func killProcess(pid int) {
	pidS := strconv.Itoa(pid)

	cmd := exec.Command(
		"sudo", "kill", "-9", pidS,
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
