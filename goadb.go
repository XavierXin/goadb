package goadb

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Device struct {
	transportID     string
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
	return d.adbCmd("shell", cmd)
}

// cmd can be "pull", "push", etc.
func (d *Device) adbCmd(cmd string, args string) (string, error) {
	splitArgs := strings.Split(args, " ")
	argsWithID := append([]string{"-t", d.transportID, cmd}, splitArgs...)
	return d.commandExecuter(d.adbPath, argsWithID...)
}

func (d *Device) IsActive() bool {
	output, err := d.ShellCmd("whoami")
	return err == nil && !strings.Contains(output, "error: no device with transport id")
}

// executeShellCmd execute "cmd args..."
// args has to include -s transportID
// isolated this function out to make testing easier
func (d *Device) executeShellCmd(cmd string, args ...string) (output string, err error) {
	execCmd := exec.Command(cmd, args...)
	out, err := execCmd.CombinedOutput()
	output = string(out)
	if len(output) != 0 {
		output = string(out[:len(out)-1])
	}
	return output, err // get rid of last \n
}

func (d *Device) HostName() (string, error) {
	return d.ShellCmd("hostname")
}

func (d *Device) Pull(src, dst string) error {
	_, err := d.adbCmd("pull", fmt.Sprintf("%s %s", src, dst))
	return err
}
