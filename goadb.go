package goadb

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type DeviceIf struct {
}

type Device struct {
	transportID     int
	adbPath         string
	commandExecuter func(cmd string, args ...string) (string, error)
}

func GetAllConnectedDevices() (devices []*Device, err error) {
	adbPath := findAdbPath()
	if len(adbPath) == 0 {
		return devices, errors.New("no installed adb is found")
	}

	return getDevices(adbPath)
}

func (d *Device) ShellCmd(cmd string) (string, error) {
	args := strings.Split(cmd, " ")
	argsWithID := append([]string{"-s", strconv.Itoa(d.transportID), "shell"}, args...)
	return d.commandExecuter(d.adbPath, argsWithID...)
}

// executeShellCmd execute "cmd args..."
// args has to include -s transportID
// isolated this function out to make testing easier
func (d *Device) executeShellCmd(cmd string, args ...string) (output string, err error) {
	execCmd := exec.Command(cmd, args...)
	out, err := execCmd.CombinedOutput()
	return string(out), err
}
