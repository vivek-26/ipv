package ipvanish

import (
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

	pidFilePath := getPIDFilepath()
	pidF, err := os.Open(pidFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			reporter.Error("Cannot find pid file")
		}
		reporter.Error(err)
	}

	pidFData, err := ioutil.ReadAll(pidF)
	if err != nil {
		reporter.Error("Cannot read pid file data")
	}

	pidS := strings.TrimSpace(string(pidFData))
	pid, err := strconv.Atoi(pidS)
	if err != nil {
		reporter.Error("Bad Process ID")
	}

	if !isProcessRunning(pid) {
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

func isProcessRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		reporter.Error(err)
	}

	// Double check if process is running and alive
	// by sending a signal 0
	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		if err == syscall.ESRCH {
			return false
		}
		reporter.Error(err)
	}

	return true
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
